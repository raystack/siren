package notification

import (
	"context"
	"time"

	"github.com/odpf/salt/log"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"gopkg.in/yaml.v3"

	"github.com/odpf/siren/core/idempotency"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
)

//go:generate mockery --name=IdempotencyRepository -r --case underscore --with-expecter --structname IdempotencyRepository --filename idempotency_repository.go --output=./mocks
type IdempotencyRepository interface {
	InsertOnConflictReturning(context.Context, string, string) (*idempotency.Idempotency, error)
	UpdateSuccess(context.Context, uint64, bool) error
	Delete(context.Context, idempotency.Filter) error
}

//go:generate mockery --name=SubscriptionService -r --case underscore --with-expecter --structname SubscriptionService --filename subscription_service.go --output=./mocks
type SubscriptionService interface {
	MatchByLabels(ctx context.Context, namespaceID uint64, labels map[string]string) ([]subscription.Subscription, error)
}

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	Get(ctx context.Context, id uint64, gopts ...receiver.GetOption) (*receiver.Receiver, error)
}

// Service is a service for notification domain
type Service struct {
	logger                log.Logger
	q                     Queuer
	idempotencyRepository IdempotencyRepository
	receiverService       ReceiverService
	subscriptionService   SubscriptionService
	notifierPlugins       map[string]Notifier
	messagingTracer       *telemetry.MessagingTracer
}

// NewService creates a new notification service
func NewService(
	logger log.Logger,
	q Queuer,
	idempotencyRepository IdempotencyRepository,
	receiverService ReceiverService,
	subscriptionService SubscriptionService,
	notifierPlugins map[string]Notifier,
) *Service {
	ns := &Service{
		logger:                logger,
		q:                     q,
		idempotencyRepository: idempotencyRepository,
		receiverService:       receiverService,
		subscriptionService:   subscriptionService,
		notifierPlugins:       notifierPlugins,
	}

	ns.messagingTracer = telemetry.NewMessagingTracer("default")
	if q != nil {
		ns.messagingTracer = telemetry.NewMessagingTracer(q.Type())
	}

	return ns
}

func (ns *Service) getNotifierPlugin(receiverType string) (Notifier, error) {
	notifierPlugin, exist := ns.notifierPlugins[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return notifierPlugin, nil
}

func (ns *Service) DispatchToReceiver(ctx context.Context, n Notification, receiverID uint64) error {
	rcv, err := ns.receiverService.Get(ctx, receiverID, receiver.GetWithData(false))
	if err != nil {
		return err
	}

	ctx, span := ns.messagingTracer.StartSpan(ctx, "prepare_enqueue",
		trace.StringAttribute("messaging.notification_id", n.ID),
		trace.StringAttribute("messaging.routing_method", RoutingMethodReceiver.String()),
	)
	defer span.End()

	notifierPlugin, err := ns.getNotifierPlugin(rcv.Type)
	if err != nil {
		return errors.ErrInvalid.WithMsgf("invalid receiver type: %s", err.Error())
	}

	message, err := n.ToMessage(rcv.Type, rcv.Configurations)
	if err != nil {
		return err
	}

	newConfigs, err := notifierPlugin.PreHookQueueTransformConfigs(ctx, message.Configs)
	if err != nil {
		telemetry.IncrementInt64Counter(ctx, telemetry.MetricReceiverHookFailed,
			tag.Upsert(telemetry.TagRoutingMethod, RoutingMethodReceiver.String()),
			tag.Upsert(telemetry.TagReceiverType, message.ReceiverType),
			tag.Upsert(telemetry.TagHookCondition, telemetry.HookConditionPreHookQueue),
		)

		return err
	}
	message.Configs = newConfigs

	message.AddStringDetail(DetailsKeyRoutingMethod, RoutingMethodReceiver.String())

	span.End()

	telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
		tag.Upsert(telemetry.TagRoutingMethod, RoutingMethodReceiver.String()),
		tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
		tag.Upsert(telemetry.TagReceiverType, message.ReceiverType),
	)

	// supported no templating for now
	if err := ns.q.Enqueue(ctx, *message); err != nil {
		return err
	}

	return nil
}

func (ns *Service) DispatchToSubscribers(ctx context.Context, namespaceID uint64, n Notification) error {
	subscriptions, err := ns.subscriptionService.MatchByLabels(ctx, namespaceID, n.Labels)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationSubscriberNotFound)
		return errors.ErrInvalid.WithMsgf("not matching any subscription")
	}

	ctx, span := ns.messagingTracer.StartSpan(ctx, "prepare_enqueue",
		trace.StringAttribute("messaging.notification_id", n.ID),
		trace.StringAttribute("messaging.routing_method", RoutingMethodSubscribers.String()),
	)
	defer span.End()

	var messages = make([]Message, 0)

	for _, s := range subscriptions {
		for _, rcv := range s.Receivers {

			notifierPlugin, err := ns.getNotifierPlugin(rcv.Type)
			if err != nil {
				return err
			}

			message, err := n.ToMessage(rcv.Type, rcv.Configuration)
			if err != nil {
				return err
			}

			newConfigs, err := notifierPlugin.PreHookQueueTransformConfigs(ctx, message.Configs)
			if err != nil {
				telemetry.IncrementInt64Counter(ctx, telemetry.MetricReceiverHookFailed,
					tag.Upsert(telemetry.TagReceiverType, message.ReceiverType),
					tag.Upsert(telemetry.TagRoutingMethod, RoutingMethodSubscribers.String()),
					tag.Upsert(telemetry.TagHookCondition, telemetry.HookConditionPreHookQueue),
				)

				return err
			}
			message.Configs = newConfigs

			message.AddStringDetail(DetailsKeyRoutingMethod, RoutingMethodSubscribers.String())

			//TODO fetch template if any, if not exist, check provider type, if exist use the default template, if not pass as-is
			// if there is template, render and replace detail with the new one
			if n.Template != "" {
				var templateBody string

				if template.IsReservedName(n.Template) {
					templateBody = notifierPlugin.GetSystemDefaultTemplate()
				}

				if templateBody != "" {
					renderedDetailString, err := template.RenderBody(templateBody, n)
					if err != nil {
						return errors.ErrInvalid.WithMsgf("failed to render template receiver %s: %s", rcv.Type, err.Error())
					}

					var messageDetails map[string]interface{}
					if err := yaml.Unmarshal([]byte(renderedDetailString), &messageDetails); err != nil {
						return errors.ErrInvalid.WithMsgf("failed to unmarshal rendered template receiver %s: %s, rendered template: %v", rcv.Type, err.Error(), renderedDetailString)
					}
					message.Details = messageDetails
				}
			}

			telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
				tag.Upsert(telemetry.TagRoutingMethod, RoutingMethodSubscribers.String()),
				tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
				tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

			messages = append(messages, *message)
		}
	}

	span.End()

	if err := ns.q.Enqueue(ctx, messages...); err != nil {
		return err
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
	return s.idempotencyRepository.Delete(ctx, idempotency.Filter{
		TTL: TTL,
	})
}
