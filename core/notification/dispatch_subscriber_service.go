package notification

import (
	"context"
	"fmt"

	saltlog "github.com/odpf/salt/log"
	"github.com/odpf/siren/core/log"
	"github.com/odpf/siren/core/silence"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
)

type DispatchSubscriberService struct {
	logger              saltlog.Logger
	subscriptionService SubscriptionService
	silenceService      SilenceService
	notifierPlugins     map[string]Notifier
}

func NewDispatchSubscriberService(
	logger saltlog.Logger,
	subscriptionService SubscriptionService,
	silenceService SilenceService,
	notifierPlugins map[string]Notifier) *DispatchSubscriberService {
	return &DispatchSubscriberService{
		logger:              logger,
		subscriptionService: subscriptionService,
		silenceService:      silenceService,
		notifierPlugins:     notifierPlugins,
	}
}

func (s *DispatchSubscriberService) getNotifierPlugin(receiverType string) (Notifier, error) {
	notifierPlugin, exist := s.notifierPlugins[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return notifierPlugin, nil
}

func (s *DispatchSubscriberService) PrepareMessage(ctx context.Context, n Notification) ([]Message, []log.Notification, bool, error) {

	var (
		messages         = make([]Message, 0)
		notificationLogs []log.Notification
		hasSilenced      bool
	)

	subscriptions, err := s.subscriptionService.MatchByLabels(ctx, n.NamespaceID, n.Labels)
	if err != nil {
		return nil, nil, false, err
	}

	if len(subscriptions) == 0 {
		telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationSubscriberNotFound)
		return nil, nil, false, errors.ErrInvalid.WithMsgf("not matching any subscription")
	}

	for _, sub := range subscriptions {

		if len(sub.Receivers) == 0 {
			s.logger.Warn(fmt.Sprintf("invalid subscription with id %d, no receiver found", sub.ID))
			continue
		}

		// try silencing by labels
		silences, err := s.silenceService.List(ctx, silence.Filter{
			NamespaceID:       n.NamespaceID,
			SubscriptionMatch: sub.Match,
		})
		if err != nil {
			return nil, nil, false, err
		}

		if len(silences) != 0 {
			hasSilenced = true

			var silenceIDs []string
			for _, sil := range silences {
				silenceIDs = append(silenceIDs, sil.ID)
			}

			notificationLogs = append(notificationLogs, log.Notification{
				NamespaceID:    n.NamespaceID,
				NotificationID: n.ID,
				SubscriptionID: sub.ID,
				AlertIDs:       n.AlertIDs,
				SilenceIDs:     silenceIDs,
			})

			s.logger.Info(fmt.Sprintf("notification '%s' of alert ids '%v' is being silenced by labels '%v'", n.ID, n.AlertIDs, silences))
			continue
		}

		// subscription not being silenced by label
		silences, err = s.silenceService.List(ctx, silence.Filter{
			NamespaceID:    n.NamespaceID,
			SubscriptionID: sub.ID,
		})
		if err != nil {
			return nil, nil, false, err
		}

		silencedReceiversMap, validReceivers, err := sub.SilenceReceivers(silences)
		if err != nil {
			return nil, nil, false, errors.ErrInvalid.WithMsgf(err.Error())
		}

		if len(silencedReceiversMap) != 0 {
			hasSilenced = true

			for rcvID, sils := range silencedReceiversMap {
				var silenceIDs []string
				for _, sil := range sils {
					silenceIDs = append(silenceIDs, sil.ID)
				}

				notificationLogs = append(notificationLogs, log.Notification{
					NamespaceID:    n.NamespaceID,
					NotificationID: n.ID,
					SubscriptionID: sub.ID,
					ReceiverID:     rcvID,
					AlertIDs:       n.AlertIDs,
					SilenceIDs:     silenceIDs,
				})
			}
		}

		for _, rcv := range validReceivers {
			notifierPlugin, err := s.getNotifierPlugin(rcv.Type)
			if err != nil {
				return nil, nil, false, err
			}

			message, err := InitMessage(
				ctx,
				notifierPlugin,
				n,
				rcv.Type,
				rcv.Configuration,
				InitWithExpiryDuration(n.ValidDuration),
			)
			if err != nil {
				return nil, nil, false, err
			}

			messages = append(messages, message)
			notificationLogs = append(notificationLogs, log.Notification{
				NamespaceID:    n.NamespaceID,
				NotificationID: n.ID,
				SubscriptionID: sub.ID,
				ReceiverID:     rcv.ID,
				AlertIDs:       n.AlertIDs,
			})
		}
	}

	return messages, notificationLogs, hasSilenced, nil
}
