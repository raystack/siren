package plugins

// Receiver plugin needs to implement this interface
//
// type ReceiverPlugin interface {
// 	PopulateDataFromConfigs(ctx context.Context, configs Configurations) (map[string]interface{}, error)
// 	Notify(ctx context.Context, configs Configurations, payloadMessage map[string]interface{}) error
// 	ValidateConfigurations(configs Configurations) error
// 	EnrichSubscriptionConfig(subsConfs map[string]string, configs Configurations) (map[string]string, error)
// 	PreHookTransformConfigs(ctx context.Context, configs Configurations) (Configurations, error)
// 	PostHookTransformConfigs(ctx context.Context, configs Configurations) (Configurations, error)
// }
