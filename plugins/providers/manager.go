package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/goto/salt/log"
	"github.com/goto/siren/plugins"
	"github.com/hashicorp/go-plugin"
)

type PluginManager struct {
	cfg           plugins.Config
	pluginClients map[string]*plugin.Client
}

func NewPluginManager(logger log.Logger, cfg plugins.Config) *PluginManager {
	return &PluginManager{
		cfg: cfg,
	}
}

func (pl *PluginManager) InitClients() map[string]*plugin.Client {
	var pluginClients = make(map[string]*plugin.Client, 0)
	for k, v := range pl.cfg.Plugins {
		// We're a host. Start by launching the plugin process.
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  v.Handshake.ProtocolVersion,
				MagicCookieKey:   v.Handshake.MagicCookieKey,
				MagicCookieValue: v.Handshake.MagicCookieValue,
			},
			Cmd:              exec.Command("sh", "-c", filepath.Join(pl.cfg.PluginPath, k)),
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			VersionedPlugins: map[int]plugin.PluginSet{
				1: {
					k: &plugins.ProviderV1beta1GRPCPlugin{},
				},
			},
		})
		pluginClients[k] = client
	}

	pl.pluginClients = pluginClients

	return pl.pluginClients
}

func (pl *PluginManager) DispenseClients(pluginClients map[string]*plugin.Client) (map[string]plugins.ProviderV1beta1, error) {
	var providerPlugins = make(map[string]plugins.ProviderV1beta1, 0)

	for k, client := range pluginClients {
		prot := client.Protocol()
		_ = prot
		rpcClient, err := client.Client()
		if err != nil {
			return nil, fmt.Errorf("error creating plugin client: %s with error %s", k, err.Error())
		} else {
			// Request the plugin
			raw, err := rpcClient.Dispense(k)
			if err != nil {
				return nil, fmt.Errorf("error dispensing plugin client: %s with error %s", k, err.Error())
			}

			providerClient := raw.(plugins.ProviderV1beta1)

			providerPlugins[k] = providerClient
		}
	}

	return providerPlugins, nil
}

func (pl *PluginManager) InitConfigs(ctx context.Context, providerPlugins map[string]plugins.ProviderV1beta1, logLevel string) error {
	for k, client := range providerPlugins {
		pluginConfig, ok := pl.cfg.Plugins[k]
		if !ok {
			return fmt.Errorf("cannot found config for provider %s", k)
		}

		var serviceConfig = make(map[string]interface{}, 0)

		if pluginConfig.ServiceConfig != nil {
			serviceConfig = pluginConfig.ServiceConfig
		}

		jsonRaw, err := json.Marshal(serviceConfig)
		if err != nil {
			return fmt.Errorf("cannot stringify config for provider %s with error %s", k, err.Error())
		}

		if err := client.SetConfig(ctx, string(jsonRaw)); err != nil {
			return fmt.Errorf("cannot set config for provider %s with error %s", k, err.Error())
		}
	}
	return nil
}

func (pl *PluginManager) Stop() {
	for _, pluginClient := range pl.pluginClients {
		if !pluginClient.Exited() {
			pluginClient.Kill()
		}
	}
}
