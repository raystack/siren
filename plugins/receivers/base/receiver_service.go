package base

import (
	"context"

	"github.com/odpf/siren/plugins"
)

// UnimplementedReceiverService is a base receiver plugin service layer for file
type UnimplementedReceiverService struct{}

func (s *UnimplementedReceiverService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *UnimplementedReceiverService) PostHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *UnimplementedReceiverService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *UnimplementedReceiverService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
	return nil, plugins.ErrNotImplemented
}
