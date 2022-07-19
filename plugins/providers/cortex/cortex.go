package cortex

type ReceiverConfig struct {
	Name           string
	Type           string
	Match          map[string]string
	Configurations map[string]string
}

type AlertManagerConfig struct {
	Receivers []ReceiverConfig
}
