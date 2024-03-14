package notification

import (
	"context"
	"fmt"
	"time"

	saltlog "github.com/goto/salt/log"

	"github.com/goto/siren/core/alert"
	"github.com/goto/siren/core/log"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/pkg/errors"
)

type Dispatcher interface {
	PrepareMessage(ctx context.Context, n Notification) ([]Message, []log.Notification, bool, error)
}

type SubscriptionService interface {
	MatchByLabels(ctx context.Context, namespaceID uint64, labels map[string]string) ([]subscription.Subscription, error)
}

type ReceiverService interface {
	List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error)
}

type SilenceService interface {
	List(ctx context.Context, filter silence.Filter) ([]silence.Silence, error)
}

type AlertService interface {
	UpdateSilenceStatus(ctx context.Context, alertIDs []int64, hasSilenced bool, hasNonSilenced bool) error
}

type LogService interface {
	LogNotifications(ctx context.Context, nlogs ...log.Notification) error
}

// Service is a service for notification domain
type Service struct {
	logger                saltlog.Logger
	cfg                   Config
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
	enableSilenceFeature  bool
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
	cfg Config,
	repository Repository,
	q Queuer,
	notifierPlugins map[string]Notifier,
	deps Deps,
	enableSilenceFeature bool,
) *Service {
	var (
		dispatchReceiverService   = deps.DispatchReceiverService
		dispatchSubscriberService = deps.DispatchSubscriberService
	)
	if deps.DispatchReceiverService == nil {
		dispatchReceiverService = NewDispatchReceiverService(deps.ReceiverService, notifierPlugins)
	}
	if deps.DispatchSubscriberService == nil {
		dispatchSubscriberService = NewDispatchSubscriberService(logger, deps.SubscriptionService, deps.SilenceService, notifierPlugins, enableSilenceFeature)
	}

	ns := &Service{
		logger:                logger,
		cfg:                   cfg,
		q:                     q,
		repository:            repository,
		idempotencyRepository: deps.IdempotencyRepository,
		logService:            deps.LogService,
		receiverService:       deps.ReceiverService,
		subscriptionService:   deps.SubscriptionService,
		silenceService:        deps.SilenceService,
		alertService:          deps.AlertService,
		dispatcher: map[string]Dispatcher{
			FlowReceiver:   dispatchReceiverService,
			FlowSubscriber: dispatchSubscriberService,
		},
		notifierPlugins:      notifierPlugins,
		enableSilenceFeature: enableSilenceFeature,
	}

	return ns
}

func (s *Service) getDispatcherFlowService(notificationFlow string) (Dispatcher, error) {
	selectedDispatcher, exist := s.dispatcher[notificationFlow]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported notification type: %q", notificationFlow)
	}
	return selectedDispatcher, nil
}

func (s *Service) Dispatch(ctx context.Context, n Notification) (string, error) {
	ctx = s.repository.WithTransaction(ctx)
	no, err := s.repository.Create(ctx, n)
	if err != nil {
		if err := s.repository.Rollback(ctx, err); err != nil {
			return "", err
		}
		return "", err
	}

	n.EnrichID(no.ID)

	switch n.Type {
	case TypeAlert:
		if err := s.dispatchAlerts(ctx, n); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return "", err
			}
			return "", err
		}
	case TypeEvent:
		if err := s.dispatchEvents(ctx, n); err != nil {
			if err := s.repository.Rollback(ctx, err); err != nil {
				return "", err
			}
			return "", err
		}
	default:
		if err := s.repository.Rollback(ctx, err); err != nil {
			return "", err
		}
		return "", errors.ErrInternal.WithMsgf("unknown notification type")
	}

	if err := s.repository.Commit(ctx); err != nil {
		return "", err
	}

	return n.ID, nil
}

func (s *Service) dispatchByFlow(ctx context.Context, n Notification, flow string) error {
	if err := n.Validate(flow); err != nil {
		return err
	}

	dispatcherService, err := s.getDispatcherFlowService(flow)
	if err != nil {
		return err
	}

	messages, notificationLogs, hasSilenced, err := dispatcherService.PrepareMessage(ctx, n)
	if err != nil {
		return err
	}

	if len(messages) == 0 && len(notificationLogs) == 0 {
		return fmt.Errorf("something wrong and no messages will be sent with notification: %v", n)
	}

	if err := s.logService.LogNotifications(ctx, notificationLogs...); err != nil {
		return fmt.Errorf("failed logging notifications: %w", err)
	}

	// Reliability of silence feature need to be tested more
	if s.enableSilenceFeature {
		if err := s.alertService.UpdateSilenceStatus(ctx, n.AlertIDs, hasSilenced, len(messages) != 0); err != nil {
			return fmt.Errorf("failed updating silence status: %w", err)
		}
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

func (s *Service) dispatchEvents(ctx context.Context, n Notification) error {
	if len(n.ReceiverSelectors) == 0 && len(n.Labels) == 0 {
		return errors.ErrInvalid.WithMsgf("no receivers found")
	}

	if len(n.ReceiverSelectors) != 0 {
		if err := s.dispatchByFlow(ctx, n, FlowReceiver); err != nil {
			return err
		}
	}

	if len(n.Labels) != 0 {
		if err := s.dispatchByFlow(ctx, n, FlowSubscriber); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) dispatchAlerts(ctx context.Context, n Notification) error {
	if err := s.dispatchByFlow(ctx, n, FlowSubscriber); err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckIdempotency(ctx context.Context, scope, key string) (string, error) {
	idempt, err := s.idempotencyRepository.Check(ctx, scope, key)
	if err != nil {
		return "", err
	}

	return idempt.NotificationID, nil
}

func (s *Service) InsertIdempotency(ctx context.Context, scope, key, notificationID string) error {
	if _, err := s.idempotencyRepository.Create(ctx, scope, key, notificationID); err != nil {
		return err
	}

	return nil
}

func (s *Service) RemoveIdempotencies(ctx context.Context, TTL time.Duration) error {
	return s.idempotencyRepository.Delete(ctx, IdempotencyFilter{
		TTL: TTL,
	})
}

// Transform alerts and populate Data and Labels to be interpolated to the system-default template
// .Data
// - id
// - status "FIRING"/"RESOLVED"
// - resource
// - template
// - metric_value
// - metric_name
// - generator_url
// - num_alerts_firing
// - dashboard
// - playbook
// - summary
// .Labels
// - severity "WARNING"/"CRITICAL"
// - alertname
// - (others labels defined in rules)
func (s *Service) BuildFromAlerts(
	alerts []alert.Alert,
	firingLen int,
	createdTime time.Time,
) ([]Notification, error) {
	if len(alerts) == 0 {
		return nil, errors.New("empty alerts")
	}

	alertsMap, err := groupByLabels(alerts, s.cfg.GroupBy)
	if err != nil {
		return nil, err
	}

	var notifications []Notification

	for hashKey, groupedAlerts := range alertsMap {
		sampleAlert := groupedAlerts[0]

		data := map[string]any{}

		mergedAnnotations := map[string][]string{}
		for _, a := range groupedAlerts {
			for k, v := range a.Annotations {
				mergedAnnotations[k] = append(mergedAnnotations[k], v)
			}
		}
		// make unique
		for k, v := range mergedAnnotations {
			mergedAnnotations[k] = removeDuplicateStringValues(v)
		}
		// render annotations
		for k, vSlice := range mergedAnnotations {
			for _, v := range vSlice {
				if _, ok := data[k]; ok {
					data[k] = fmt.Sprintf("%s\n%s", data[k], v)
				} else {
					data[k] = v
				}
			}
		}

		data["status"] = sampleAlert.Status
		data["generator_url"] = sampleAlert.GeneratorURL
		data["num_alerts_firing"] = firingLen

		alertIDs := []int64{}

		for _, a := range groupedAlerts {
			alertIDs = append(alertIDs, int64(a.ID))
		}

		for k, v := range sampleAlert.Labels {
			data[k] = v
		}

		notifications = append(notifications, Notification{
			NamespaceID: sampleAlert.NamespaceID,
			Type:        TypeAlert,
			Data:        data,
			Labels:      sampleAlert.Labels,
			Template:    template.ReservedName_SystemDefault,
			UniqueKey:   hashGroupKey(sampleAlert.GroupKey, hashKey),
			CreatedAt:   createdTime,
			AlertIDs:    alertIDs,
		})
	}

	return notifications, nil
}
