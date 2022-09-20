package receiver

import "context"

// Resolver is an interface for the receiver to resolve all configs and function
// related to a specific receiver type. Receiver plugin needs to implement
// this interface for all configs and functionality resolution.
//
//go:generate mockery --name=Resolver -r --case underscore --with-expecter --structname Resolver --filename resolver.go --output=./mocks
type Resolver interface {
	PopulateDataFromConfigs(ctx context.Context, configs Configurations) (map[string]interface{}, error)
	Notify(ctx context.Context, configs Configurations, payloadMessage map[string]interface{}) error
	ValidateConfigurations(configs Configurations) error
	EnrichSubscriptionConfig(subsConfs map[string]string, configs Configurations) (map[string]string, error)
	PreHookTransformConfigs(ctx context.Context, configs Configurations) (Configurations, error)
	PostHookTransformConfigs(ctx context.Context, configs Configurations) (Configurations, error)
}
