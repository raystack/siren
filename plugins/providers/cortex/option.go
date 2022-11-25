package cortex

import "github.com/odpf/siren/pkg/httpclient"

type ServiceOption func(*PluginService)

// WithCortexClient uses cortex-tools client passed in the argument
func WithCortexClient(cc CortexCaller) ServiceOption {
	return func(so *PluginService) {
		so.cortexClient = cc
	}
}

// WithHTTPClient assigns custom client when creating a http client
func WithHTTPClient(cli *httpclient.Client) ServiceOption {
	return func(so *PluginService) {
		so.httpClient = cli
	}
}
