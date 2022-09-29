package base

import (
	"context"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/plugins"
)

// UnimplementedNotificationService is a base notification plugin service layer for File webhook
type UnimplementedNotificationService struct{}

func (s *UnimplementedNotificationService) ValidateConfigMap(notificationConfigMap map[string]interface{}) error {
	return plugins.ErrNotImplemented
}

func (s *UnimplementedNotificationService) Publish(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	return false, plugins.ErrNotImplemented
}

func (s *UnimplementedNotificationService) PreHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	return notificationConfigMap, nil
}

func (s *UnimplementedNotificationService) PostHookTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	return notificationConfigMap, nil
}

func (s *UnimplementedNotificationService) DefaultTemplateOfProvider(providerType string) string {
	switch providerType {
	case provider.TypeCortex:
		return ""
	default:
		return ""
	}
}
