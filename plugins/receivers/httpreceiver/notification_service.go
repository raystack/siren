package httpreceiver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/plugins/receivers/base"
)

// NotificationService is a notification plugin service layer for http webhook
type NotificationService struct {
	base.UnimplementedNotificationService
	client HTTPCaller
}

// NewNotificationService returns httpreceiver service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(client HTTPCaller) *NotificationService {
	return &NotificationService{
		client: client,
	}
}

func (h *NotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	bodyBytes, err := json.Marshal(notificationMessage.Details)
	if err != nil {
		return false, err
	}

	if err := h.client.Notify(ctx, notificationConfig.URL, bodyBytes); err != nil {
		if errors.As(err, new(retry.RetryableError)) {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (h *NotificationService) PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to slack notification config: %w", err)
	}

	if err := notificationConfig.Validate(); err != nil {
		return nil, err
	}

	return notificationConfig.AsMap(), nil
}
