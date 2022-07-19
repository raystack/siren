package pagerduty

import (
	"context"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins"
)

type PagerDutyService struct{}

// NewReceiverService returns pagerduty service struct
func NewReceiverService() *PagerDutyService {
	return &PagerDutyService{}
}

func (s *PagerDutyService) Notify(ctx context.Context, configurations receiver.Configurations, payloadMessage map[string]interface{}) error {
	return plugins.ErrNotImplemented
}

func (s *PagerDutyService) PreHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	return configurations, nil
}

func (s *PagerDutyService) PostHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	return configurations, nil
}

func (s *PagerDutyService) PopulateDataFromConfigs(ctx context.Context, configurations receiver.Configurations) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *PagerDutyService) ValidateConfigurations(configurations receiver.Configurations) error {
	_, err := configurations.GetString("service_key")
	if err != nil {
		return err
	}

	return nil
}

func (s *PagerDutyService) EnrichSubscriptionConfig(subsConfs map[string]string, receiverConfs receiver.Configurations) (map[string]string, error) {
	mapConf := make(map[string]string)
	if val, ok := receiverConfs["service_key"]; ok {
		if mapConf["service_key"], ok = val.(string); !ok {
			return nil, errors.New("service_key config from receiver should be in string")
		}
	}
	return mapConf, nil
}
