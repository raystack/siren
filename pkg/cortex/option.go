package cortex

type ClientOption func(*Client)

func WithHelperTemplate(configYaml, helperTemplate string) ClientOption {
	return func(c *Client) {
		c.configYaml = configYaml
		c.helperTemplate = helperTemplate
	}
}

func WithCortexClient(cc CortexCaller) ClientOption {
	return func(c *Client) {
		c.cortexClient = cc
	}
}
