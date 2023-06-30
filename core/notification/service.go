package notification

import (
	"context"
	"fmt"
	"time"

	saltlog "github.com/raystack/salt/log"
	"go.opencensus.io/trace"

	"github.com/raystack/siren/core/log"
	"github.com/raystack/siren/core/receiver"
	"github.com/raystack/siren/core/silence"
	"github.com/raystack/siren/core/subscription"
	"github.com/raystack/siren/pkg/errors"
	"github.com/raystack/siren/pkg/telemetry"
)

//go:generate mockery --name=Dispatcher -r --case underscore --with-expecter --structname Dispatcher --filename dispatcher.go --output=./mocks
type Dispatcher interface {
	PrepareMessage(ctx context.Context, n Notification) ([]Message, []log.Notification, bool, error)
}

//go:generate mockery --name=SubscriptionService -r --case underscore --with-expecter --structname SubscriptionService --filename subscription_service.go --output=./mocks
type SubscriptionService interface {
	MatchByLabels(ctx context.Context, namespaceID uint64, labels map[string]string) ([]subscription.Subscription, error)
}

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	Get(ctx context.Context, id uint64, gopts ...receiver.GetOption) (*receiver.Receiver, error)
}

//go:generate mockery --name=SilenceService -r --case underscore --with-expecter --structname SilenceService --filename silence_service.go --output=./mocks
type SilenceService interface {
	List(ctx context.Context, filter silence.Filter) ([]silence.Silence, error)
}

//go:generate mockery --name=AlertService -r --case underscore --with-expecter --structname AlertService --filename alert_service.go --output=./mocks
type AlertService interface {
	UpdateSilenceStatus(ctx context.Context, alertIDs []int64, hasSilenced bool, hasNonSilenced bool) error
}

//go:generate mockery --name=LogService -r --case underscore --with-expecter --structname LogService --filename log_service.go --output=./mocks
type LogService interface {
	LogNotifications(ctx context.Context, nlogs ...log.Notification) error
}

// Service is a service for notification domain
type Service struct {
	logger                saltlog.Logger
	q                     Queuer
	idempotencyRepository IdempotencyRepository
	logService            LogService
	repository            Repository
	receiverService       ReceiverService
	subscriptionService   SubscriptionService
	silenceService        SilenceService
	alertService          AlertService
	notifierPlugins       map[string]Notifier
	dispatcher            map[string]Dispatcher
	messagingTracer       *telemetry.MessagingTracer
}

type Deps struct {
	IdempotencyRepository     IdempotencyRepository
	LogService                LogService
	ReceiverService           ReceiverService
	SubscriptionService       SubscriptionService
	SilenceService            SilenceService
	AlertService              AlertService
	DispatchReceiverService   Dispatcher
	DispatchSubscriberService Dispatcher
}

// NewService creates a new notification service
func NewService(
	logger saltlog.Logger,
	repository Repository,
	q Queuer,
	notifierPlugins map[string]Notifier,
	deps Deps,
) *Service {
	var (
		dispatchReceiverService   = deps.DispatchReceiverService
		dispatchSubscriberService = deps.DispatchSubscriberService
	)
	if deps.DispatchReceiverService == nil {
		dispatchReceiverService = NewDispatchReceiverService(deps.ReceiverService, notifierPlugins)
	}
	if deps.DispatchSubscriberService == nil {
		dispatchSubscriberService = NewDispatchSubscriberService(logger, deps.SubscriptionService, deps.SilenceService, notifierPlugins)
	}

	ns := &Service{
		logger:                logger,
		q:                     q,
		repository:            repository,
		idempotencyRepository: deps.IdempotencyRepository,
		logService:            deps.LogService,
		receiverService:       deps.ReceiverService,
		subscriptionService:   deps.SubscriptionService,
		silenceService:        deps.SilenceService,
		alertService:          deps.AlertService,
		dispatcher: map[string]Dispatcher{
			TypeReceiver:   dispatchReceiverService,
			TypeSubscriber: dispatchSubscriberService,
		},
		notifierPlugins: notifierPlugins,
	}

	ns.messagingTracer = telemetry.NewMessagingTracer("default")
	if q != nil {
		ns.messagingTracer = telemetry.NewMessagingTracer(q.Type())
	}

	return ns
}

func (s *Service) getDispatcherService(notificationType string) (Dispatcher, error) {
	selectedDispatcher, exist := s.dispatcher[notificationType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported notification type: %q", notificationType)
	}
	return selectedDispatcher, nil
}

func (s *Service) Dispatch(ctx context.Context, n Notification) error {
	if err := n.Validate(); err != nil {
		return err
	}

	no, err := s.repository.Create(ctx, n)
	if err != nil {
		return err
	}

	n.EnrichID(no.ID)

	dispatcherService, err := s.getDispatcherService(n.Type)
	if err != nil {
		return err
	}

	ctx, span := s.messagingTracer.StartSpan(ctx, "prepare_message",
		trace.StringAttribute("messaging.notification_id", n.ID),
		trace.StringAttribute("messaging.routing_method", n.Type),
	)
	messages, notificationLogs, hasSilenced, err := dispatcherService.PrepareMessage(ctx, n)
	span.End()
	if err != nil {
		return err
	}

	if len(messages) == 0 && len(notificationLogs) == 0 {
		return fmt.Errorf("something wrong and no messages will be sent with notification: %v", n)
	}

	if err := s.logService.LogNotifications(ctx, notificationLogs...); err != nil {
		return fmt.Errorf("failed logging notifications: %w", err)
	}

	if err := s.alertService.UpdateSilenceStatus(ctx, n.AlertIDs, hasSilenced, len(messages) != 0); err != nil {
		return fmt.Errorf("failed updating silence status: %w", err)
	}

	if len(messages) == 0 {
		s.logger.Info("no messages to enqueue")
		return nil
	}

	if err := s.q.Enqueue(ctx, messages...); err != nil {
		return fmt.Errorf("failed enqueuing messages: %w", err)
	}

	return nil
}

func (s *Service) CheckAndInsertIdempotency(ctx context.Context, scope, key string) (uint64, error) {
	idempt, err := s.idempotencyRepository.InsertOnConflictReturning(ctx, scope, key)
	if err != nil {
		return 0, err
	}

	if idempt.Success {
		return 0, errors.ErrConflict
	}

	return idempt.ID, nil
}

func (s *Service) MarkIdempotencyAsSuccess(ctx context.Context, id uint64) error {
	return s.idempotencyRepository.UpdateSuccess(ctx, id, true)
}

func (s *Service) RemoveIdempotencies(ctx context.Context, TTL time.Duration) error {
	return s.idempotencyRepository.Delete(ctx, IdempotencyFilter{
		TTL: TTL,
	})
}
