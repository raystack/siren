package httpreceiver

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/retry"
)

// HTTPNotificationService is a notification plugin service layer for http webhook
type HTTPNotificationService struct {
	client HTTPCaller
}

// NewNotificationService returns httpreceiver service struct. This service implement [receiver.Notifier] interface.
func NewNotificationService(client HTTPCaller) *HTTPNotificationService {
	return &HTTPNotificationService{
		client: client,
	}
}

func (h *HTTPNotificationService) ValidateConfigMap(notificationConfigMap map[string]interface{}) error {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationConfigMap, notificationConfig); err != nil {
		return err
	}

	if err := notificationConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (h *HTTPNotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	notificationConfig := &NotificationConfig{}
	if err := mapstructure.Decode(notificationMessage.Configs, notificationConfig); err != nil {
		return false, err
	}

	bodyBytes, err := json.Marshal(notificationMessage.Detail)
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
