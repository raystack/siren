package receiver

import "context"

// ConfigResolver is an interface for the receiver to resolve all configs and function
// related to a specific receiver type. Receiver plugin needs to implement
// this interface for all configs and functionality resolution.
//

type ConfigResolver interface {
	BuildData(ctx context.Context, configs map[string]any) (map[string]any, error)
	PreHookDBTransformConfigs(ctx context.Context, configs map[string]any) (map[string]any, error)
	PostHookDBTransformConfigs(ctx context.Context, configs map[string]any) (map[string]any, error)
}
