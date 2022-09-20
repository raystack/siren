package httpreceiver

import (
	"context"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins"
)

// HTTPService is a receiver plugin service layer for http
type HTTPService struct{}

// NewReceiverService returns httpreceiver service struct. This service implement [receiver.Resolver] interface.
func NewReceiverService() *HTTPService {
	return &HTTPService{}
}

func (s *HTTPService) Notify(ctx context.Context, configurations receiver.Configurations, payloadMessage map[string]interface{}) error {
	return plugins.ErrNotImplemented
}

func (s *HTTPService) PreHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	return configurations, nil
}

func (s *HTTPService) PostHookTransformConfigs(ctx context.Context, configurations receiver.Configurations) (receiver.Configurations, error) {
	return configurations, nil
}

func (s *HTTPService) PopulateDataFromConfigs(ctx context.Context, configurations receiver.Configurations) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *HTTPService) ValidateConfigurations(configurations receiver.Configurations) error {
	_, err := configurations.GetString("url")
	if err != nil {
		return err
	}
	return nil
}

func (s *HTTPService) EnrichSubscriptionConfig(subsConfs map[string]string, receiverConfs receiver.Configurations) (map[string]string, error) {
	mapConf := make(map[string]string)
	if val, ok := receiverConfs["url"]; ok {
		if mapConf["url"], ok = val.(string); !ok {
			return nil, errors.New("url config from receiver should be in string")
		}
	}
	return mapConf, nil
}
