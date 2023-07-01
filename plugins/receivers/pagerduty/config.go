package pagerduty

import (
	"errors"
	"fmt"

	"github.com/raystack/siren/pkg/httpclient"
	"github.com/raystack/siren/pkg/retry"
	"github.com/raystack/siren/pkg/secret"
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

func (c *ReceiverConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"service_key": c.ServiceKey,
	}
}

type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
}

func (c *NotificationConfig) AsMap() map[string]interface{} {
	return c.ReceiverConfig.AsMap()
}
