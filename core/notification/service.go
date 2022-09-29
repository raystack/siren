package notification

import (
	"context"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/errors"
)

type NotificationService struct {
	logger              log.Logger
	q                   Queuer
	subscriptionService SubscriptionService
	receiverPlugins     map[string]Notifier
}

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

			if err := receiverPlugin.ValidateConfig(message.Configs); err != nil {
				return err
			}

			messages = append(messages, *message)
		}
	}

	if err := ns.q.Enqueue(ctx, messages...); err != nil {
		return err
	}

	return nil
}
