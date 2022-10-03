package pagerduty

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins"
)

// PagerDutyService is a receiver plugin service layer for pagerduty
type PagerDutyService struct{}

// NewService returns pagerduty service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService() *PagerDutyService {
	return &PagerDutyService{}
}

func (s *PagerDutyService) Notify(ctx context.Context, configurations map[string]interface{}, payloadMessage map[string]interface{}) error {
	return plugins.ErrNotImplemented
}

func (s *PagerDutyService) PreHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(configurations, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	if err := receiverConfig.Validate(); err != nil {
		return nil, errors.ErrInvalid.WithMsgf(err.Error())
	}

	return configurations, nil
}

func (s *PagerDutyService) PostHookTransformConfigs(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return configurations, nil
}

func (s *PagerDutyService) BuildData(ctx context.Context, configurations map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *PagerDutyService) BuildNotificationConfig(subsConfs map[string]interface{}, receiverConfs map[string]interface{}) (map[string]interface{}, error) {
	receiverConfig := &ReceiverConfig{}
	if err := mapstructure.Decode(receiverConfs, receiverConfig); err != nil {
		return nil, fmt.Errorf("failed to transform configurations to receiver config: %w", err)
	}

	notificationConfig := NotificationConfig{
		ReceiverConfig: *receiverConfig,
	}

	return notificationConfig.AsMap(), nil
}
