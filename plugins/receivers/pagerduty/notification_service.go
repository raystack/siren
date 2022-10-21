package pagerduty

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/plugins/receivers/base"
)

// NotificationService is a notification plugin service layer for pagerduty
type NotificationService struct {
	base.UnimplementedNotificationService
	client PagerDutyCaller
}

// NewNotificationService returns pagerduty service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(client PagerDutyCaller) *NotificationService {
	return &NotificationService{
		client: client,
	}
}

func (pd *NotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	pgMessageV1 := &MessageV1{}
	if err := mapstructure.Decode(notificationMessage.Details, &pgMessageV1); err != nil {
		return false, err
	}
	pgMessageV1.ServiceKey = notificationConfig.ServiceKey

	if err := pd.client.NotifyV1(ctx, *pgMessageV1); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (pd *NotificationService) PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to slack notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	return notificationConfig.AsMap(), nil
}

func (pd *NotificationService) DefaultTemplateOfProvider(templateName string) string {
	switch templateName {
	case template.ReservedName_DefaultCortex:
		return defaultCortexAlertTemplateBodyV1
	default:
		return ""
	}
}
