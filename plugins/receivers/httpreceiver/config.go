package httpreceiver

import (
	"fmt"

	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/pkg/retry"
)

// AppConfig is a config loaded when siren is started
type AppConfig struct {
	Retry      retry.Config      `mapstructure:"retry" yaml:"retry"`
	HTTPClient httpclient.Config `mapstructure:"http_client" yaml:"http_client"`
}
type ReceiverConfig struct {
	URL string `mapstructure:"url"`
}

func (c *ReceiverConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("invalid http receiver config, url: %s", c.URL)
	}
	return nil
}

func (c *ReceiverConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"url": c.URL,
	}
}

type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
}

func (c *NotificationConfig) AsMap() map[string]interface{} {
	return c.ReceiverConfig.AsMap()
}
