package httpreceiver

import "fmt"

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
	ReceiverConfig
}

func (c *NotificationConfig) AsMap() map[string]interface{} {
	return c.ReceiverConfig.AsMap()
}
