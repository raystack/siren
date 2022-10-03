package receiver

import "context"

// Resolver is an interface for the receiver to resolve all configs and function
// related to a specific receiver type. Receiver plugin needs to implement
// this interface for all configs and functionality resolution.
//
//go:generate mockery --name=Resolver -r --case underscore --with-expecter --structname Resolver --filename resolver.go --output=./mocks
type Resolver interface {
	BuildData(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
	BuildNotificationConfig(subscriptionConfigMap map[string]interface{}, receiverConfigMap map[string]interface{}) (map[string]interface{}, error)
	PreHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
	PostHookTransformConfigs(ctx context.Context, configs map[string]interface{}) (map[string]interface{}, error)
}
