package base

import (
	"context"

	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/plugins"
)

// UnimplementedService is a base receiver plugin service layer
type UnimplementedService struct{}

func (s *UnimplementedService) PreHookDBTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *UnimplementedService) PostHookDBTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *UnimplementedService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *UnimplementedService) PreHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	return notificationConfigMap, nil
}

func (s *UnimplementedService) PostHookQueueTransformConfigs(ctx context.Context, notificationConfigMap map[string]interface{}) (map[string]interface{}, error) {
	return notificationConfigMap, nil
}

func (s *UnimplementedService) GetSystemDefaultTemplate() string {
	return ""
}

func (s *UnimplementedService) Send(ctx context.Context, notificationMessage notification.Message) (bool, error) {
	return false, plugins.ErrNotImplemented
}
