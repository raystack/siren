package cortex

type ServiceOption func(*PluginService)

// WithCortexClient uses cortex-tools client passed in the argument
func WithCortexClient(cc CortexCaller) ServiceOption {
	return func(so *PluginService) {
		so.cortexClient = cc
	}
}
