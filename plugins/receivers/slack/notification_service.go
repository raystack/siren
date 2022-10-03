package slack

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
)

// SlackNotificationService is a notification plugin service layer for slack
type SlackNotificationService struct {
	slackClient SlackClient
}

// NewNotificationService returns slack service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(slackClient SlackClient) *SlackNotificationService {
	return &SlackNotificationService{
		slackClient: slackClient,
	}
}

func (s *SlackNotificationService) ValidateConfig(notificationConfigMap map[string]interface{}) error {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return err
	}

	if err := notificationConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *SlackNotificationService) Publish(ctx context.Context, notificationMessage notification.Message) error {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return err
	}

	slackMessage := &Message{}
	if err := mapstructure.Decode(notificationMessage.Detail, &slackMessage); err != nil {
		return err
	}

	if notificationConfig.ChannelType == "" {
		notificationConfig.ChannelType = DefaultChannelType
	}

	if err := s.slackClient.Notify(ctx, *notificationConfig, *slackMessage, CallWithToken(notificationConfig.Token)); err != nil {
		return fmt.Errorf("error calling slack notify: %w", err)
	}

	return nil
}
