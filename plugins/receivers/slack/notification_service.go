package slack

import (
	"context"
	"encoding/json"
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
	goslackMessage := &MessageGoSlack{}
	if err := goslackMessage.FromNotificationMessage(notificationMessage); err != nil {
		return err
	}

	notificationConfig := &NotificationConfig{}
	if err := notificationConfig.FromNotificationMessage(notificationMessage); err != nil {
		return err
	}

	if err := s.slackClient.Notify(ctx, goslackMessage, CallWithToken(notificationConfig.Token)); err != nil {
		return fmt.Errorf("error calling slack notify: %w", err)
	}

	return nil
}

// ToSlackMessage
//
//	{
//		"receiver_name": "",
//		"receiver_type": "",
//		"message": "",
//		"blocks": [
//				{
//					"": ""
//				}
//			]
//	}
func GetSlackMessage(payloadMessage map[string]interface{}) (*MessageGoSlack, error) {
	jsonByte, err := json.Marshal(payloadMessage)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal notification message: %w", err)
	}

	sm := &MessageGoSlack{}
	if err := json.Unmarshal(jsonByte, sm); err != nil {
		return nil, fmt.Errorf("unable to unmarshal notification message byte to slack message: %w", err)
	}

	if err := sm.Validate(); err != nil {
		return nil, err
	}

	return sm, nil
}
