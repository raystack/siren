package pagerduty

import (
	"context"
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/pkg/retry"
)

// PagerDutyNotificationService is a notification plugin service layer for pagerduty
type PagerDutyNotificationService struct {
	client PagerDutyCaller
}

// NewNotificationService returns pagerduty service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(client PagerDutyCaller) *PagerDutyNotificationService {
	return &PagerDutyNotificationService{
		client: client,
	}
}

func (pd *PagerDutyNotificationService) ValidateConfigMap(notificationConfigMap map[string]interface{}) error {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return err
	}

	if err := notificationConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (pd *PagerDutyNotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
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

func (pd *PagerDutyNotificationService) DefaultTemplateOfProvider(providerType string) string {
	switch providerType {
	case provider.TypeCortex:
		return defaultCortexAlertTemplateBodyV1
	default:
		return ""
	}
}
