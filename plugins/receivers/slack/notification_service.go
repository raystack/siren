package slack

import (
	"context"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/retry"
)

const (
	TypeChannelChannel = "channel"
	TypeChannelUser    = "user"

	defaultChannelType = TypeChannelChannel
)

// SlackNotificationService is a notification plugin service layer for slack
type SlackNotificationService struct {
	client SlackCaller
}

// NewNotificationService returns slack service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(client SlackCaller) *SlackNotificationService {
	return &SlackNotificationService{
		client: client,
	}
}

func (s *SlackNotificationService) ValidateConfigMap(notificationConfigMap map[string]interface{}) error {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return err
	}

	if err := notificationConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (s *SlackNotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	slackMessage := &Message{}
	if err := mapstructure.Decode(notificationMessage.Details, &slackMessage); err != nil {
		return false, err
	}

	if notificationConfig.ChannelType == "" {
		notificationConfig.ChannelType = defaultChannelType
	}
	if notificationConfig.ChannelName != "" {
		slackMessage.Channel = notificationConfig.ChannelName
	}

	if err := s.client.Notify(ctx, *notificationConfig, *slackMessage); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (s *SlackNotificationService) DefaultTemplateOfProvider(templateName string) string {
	switch templateName {
	case template.ReservedName_DefaultCortex:
		return defaultCortexAlertTemplateBody
	default:
		return ""
	}
}
