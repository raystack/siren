package main

import (
	"os"

	"github.com/goto/siren/plugins"
	"github.com/goto/siren/plugins/providers/cortex/common"
	cortexv1plugin "github.com/goto/siren/plugins/providers/cortex/v1"
	sirenproviderv1beta1 "github.com/goto/siren/proto/gotocompany/siren/provider/v1beta1"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func main() {
	logLevel := common.DefaultLogLevel
	if envLogLevel := os.Getenv(common.EnvKeyLogLevel); envLogLevel != "" {
		logLevel = envLogLevel
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  common.ServiceName,
		Level: hclog.LevelFromString(logLevel),
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: cortexv1plugin.Handshake,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: map[string]plugin.Plugin{
				common.PluginName: &plugins.ProviderV1beta1GRPCPlugin{
					GRPCProvider: func() sirenproviderv1beta1.ProviderServiceServer {
						return plugins.NewProviderServer(cortexv1plugin.NewPluginService(logger))
					},
				},
			},
		},
		Logger: logger,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
