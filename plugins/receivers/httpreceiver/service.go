package httpreceiver

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins"
)

// HTTPReceiverService is a receiver plugin service layer for http
type HTTPReceiverService struct{}

// NewReceiverService returns httpreceiver service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService() *HTTPReceiverService {
	return &HTTPReceiverService{}
}

func (s *HTTPReceiverService) Notify(ctx context.Context, configurations map[string]interface{}, payloadMessage map[string]interface{}) error {
	return plugins.ErrNotImplemented
}

func (s *HTTPReceiverService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return configurations, nil
}

func (s *HTTPReceiverService) PostHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *HTTPReceiverService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *HTTPReceiverService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(receiverConfig, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	notificationConfig := NotificationConfig{
		ReceiverConfig: *receiverConfig,
	}

	return notificationConfig.AsMap(), nil
}
