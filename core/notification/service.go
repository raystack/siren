package notification

import (
	"context"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
	"gopkg.in/yaml.v3"
)

//go:generate mockery --name=SubscriptionService -r --case underscore --with-expecter --structname SubscriptionService --filename subscription_service.go --output=./mocks
type SubscriptionService interface {
	MatchByLabels(ctx context.Context, labels map[string]string) ([]subscription.Subscription, error)
}

//go:generate mockery --name=ReceiverService -r --case underscore --with-expecter --structname ReceiverService --filename receiver_service.go --output=./mocks
type ReceiverService interface {
	Get(ctx context.Context, id uint64) (*receiver.Receiver, error)
}

// NotificationService is a service for notification domain
type NotificationService struct {
	logger              log.Logger
	q                   Queuer
	receiverService     ReceiverService
	subscriptionService SubscriptionService
	notifierPlugins     map[string]Notifier
	messagingTracer     *telemetry.MessagingSpan
}

// NewService creates a new notification service
func NewService(
	logger log.Logger,
	q Queuer,
	receiverService ReceiverService,
	subscriptionService SubscriptionService,
	notifierPlugins map[string]Notifier,
) *NotificationService {
	return &NotificationService{
		logger:              logger,
		q:                   q,
		receiverService:     receiverService,
		subscriptionService: subscriptionService,
		notifierPlugins:     notifierPlugins,
		messagingTracer:     telemetry.InitMessagingSpan(q.Type(), "messages"),
	}
}

func (ns *NotificationService) getNotifierPlugin(receiverType string) (Notifier, error) {
	notifierPlugin, exist := ns.notifierPlugins[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return notifierPlugin, nil
}

func (ns *NotificationService) DispatchToReceiver(ctx context.Context, n Notification, receiverID uint64) error {
	rcv, err := ns.receiverService.Get(ctx, receiverID)
	if err != nil {
		return err
	}

	notifierPlugin, err := ns.getNotifierPlugin(rcv.Type)
	if err != nil {
		return err
	}

	message, err := n.ToMessage(ctx, rcv.Type, rcv.Configurations)
	if err != nil {
		return err
	}

	newConfigs, err := notifierPlugin.PreHookQueueTransformConfigs(ctx, message.Configs)
	if err != nil {
		return err
	}
	message.Configs = newConfigs

	message.AddStringDetail(DetailsKeyRoutingMethod, RoutingMethodReceiver.String())

	ctx, span := ns.messagingTracer.StartSpan(ctx, "prepare_enqueue", message.ID, map[string]string{
		"messages.notification_id": n.ID,
		"messages.routing_method":  string(RoutingMethodReceiver),
	})
	span.End()

	// supported no templating for now
	if err := ns.q.Enqueue(ctx, *message); err != nil {
		return err
	}

	return nil
}

func (ns *NotificationService) DispatchToSubscribers(ctx context.Context, n Notification) error {
	subscriptions, err := ns.subscriptionService.MatchByLabels(ctx, n.Labels)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return errors.ErrInvalid.WithMsgf("not matching any subscription")
	}

	var messages = make([]Message, 0)

	for _, s := range subscriptions {
		for _, rcv := range s.Receivers {

			notifierPlugin, err := ns.getNotifierPlugin(rcv.Type)
			if err != nil {
				return err
			}

			message, err := n.ToMessage(ctx, rcv.Type, rcv.Configuration)
			if err != nil {
				return err
			}

			ctx, span := ns.messagingTracer.StartSpan(ctx, "prepare_enqueue", message.ID, map[string]string{
				"messages.notification_id": n.ID,
				"messages.routing_method":  string(RoutingMethodSubscribers),
			})
			defer span.End()

			newConfigs, err := notifierPlugin.PreHookQueueTransformConfigs(ctx, message.Configs)
			if err != nil {
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
						return errors.ErrInvalid.WithMsgf(err.Error())
					}

					var messageDetails map[string]interface{}
					if err := yaml.Unmarshal([]byte(renderedDetailString), &messageDetails); err != nil {
						return err
					}
					message.Details = messageDetails
				}
			}
			span.End()

			messages = append(messages, *message)
		}
	}

	ctx, span := ns.messagingTracer.StartSpan(ctx, "batch_enqueue", "", map[string]string{
		"messages.notification_id": n.ID,
	})
	defer span.End()

	if err := ns.q.Enqueue(ctx, messages...); err != nil {
		return err
	}

	return nil
}
