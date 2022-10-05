package notification

import (
	"context"
	"fmt"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"gopkg.in/yaml.v3"
)

// NotificationService is a service for notification domain
type NotificationService struct {
	logger              log.Logger
	q                   Queuer
	subscriptionService SubscriptionService
	receiverPlugins     map[string]Notifier
}

// NewService creates a new notification service
func NewService(
	logger log.Logger,
	q Queuer,
	subscriptionService SubscriptionService,
	receiverPlugins map[string]Notifier,
) *NotificationService {
	return &NotificationService{
		logger:              logger,
		q:                   q,
		subscriptionService: subscriptionService,
		receiverPlugins:     receiverPlugins,
	}
}

func (ns *NotificationService) getReceiverPlugin(receiverType string) (Notifier, error) {
	receiverPlugin, exist := ns.receiverPlugins[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return receiverPlugin, nil
}

// Dispatch sends notification to the specified queue
func (ns *NotificationService) Dispatch(ctx context.Context, n Notification) error {
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
			message, err := n.ToMessage(rcv.Type, rcv.Configuration)
			if err != nil {
				return err
			}

			receiverPlugin, err := ns.getReceiverPlugin(rcv.Type)
			if err != nil {
				return err
			}

			if err := receiverPlugin.ValidateConfigMap(message.Configs); err != nil {
				return err
			}

			//TODO fetch template if any, if not exist, check provider type, if exist use the default template, if not pass as-is
			// if there is template, render and replace detail with the new one
			if n.ProviderType != "" {
				if templateBody := receiverPlugin.DefaultTemplateOfProvider(n.ProviderType); templateBody != "" {
					renderedDetailString, err := template.RenderBody(templateBody, n)
					if err != nil {
						return err // TODO add more error detail
					}

					// renderedDetailString = strings.TrimSpace(renderedDetailString)

					var messageDetail map[string]interface{}
					fmt.Println(renderedDetailString)
					if err := yaml.Unmarshal([]byte(renderedDetailString), &messageDetail); err != nil {
						return err // TODO add more error detail
					}
					message.Detail = messageDetail
				}
			}

			messages = append(messages, *message)
		}
	}

	if err := ns.q.Enqueue(ctx, messages...); err != nil {
		return err
	}

	return nil
}
