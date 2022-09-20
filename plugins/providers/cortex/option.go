package cortex

type ClientOption func(*Client)

// WithHelperTemplate assigns helper template and config yaml string
func WithHelperTemplate(configYaml, helperTemplate string) ClientOption {
	return func(c *Client) {
		c.configYaml = configYaml
		c.helperTemplate = helperTemplate
	}
}

// WithCortexClient uses cortex-tools client passed in the argument
func WithCortexClient(cc CortexCaller) ClientOption {
	return func(c *Client) {
		c.cortexClient = cc
	}
}
