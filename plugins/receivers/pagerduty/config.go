package pagerduty

import (
	"errors"
	"fmt"

	"github.com/goto/siren/pkg/httpclient"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/pkg/secret"
)

// AppConfig is a config loaded when siren is started
type AppConfig struct {
	APIHost    string            `mapstructure:"api_host" yaml:"api_host"`
	Retry      retry.Config      `mapstructure:"retry" yaml:"retry"`
	HTTPClient httpclient.Config `mapstructure:"http_client" yaml:"http_client"`
}

func (c AppConfig) Validate() error {
	if c.APIHost == "" {
		return errors.New("invalid pagerduty app config")
	}
	return nil
}

// TODO need to support versioning later v1 and v2
type ReceiverConfig struct {
	ServiceKey secret.MaskableString `mapstructure:"service_key"`
}

func (c *ReceiverConfig) Validate() error {
	if c.ServiceKey == "" {
		return fmt.Errorf("invalid pagerduty receiver config, service_key: %s", c.ServiceKey)
	}
	return nil
}

func (c *ReceiverConfig) AsMap() map[string]any {
	return map[string]any{
		"service_key": c.ServiceKey,
	}
}

type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
}

func (c *NotificationConfig) AsMap() map[string]any {
	return c.ReceiverConfig.AsMap()
}
