package cortex

// ReceiverConfig is a receiver configuration in cortex alertmanager format
type ReceiverConfig struct {
	Name           string
	Type           string
	Match          map[string]string
	Configurations map[string]string
}

// AlertManagerConfig is a placeholder to store cortex alertmanager format
type AlertManagerConfig struct {
	Receivers []ReceiverConfig
}
