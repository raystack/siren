package file

import (
	"fmt"
)

type ReceiverConfig struct {
	URL string `mapstructure:"url"`
}

func (c *ReceiverConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("invalid file receiver config, url: %s", c.URL)
	}
	return nil
}

func (c *ReceiverConfig) AsMap() map[string]any {
	return map[string]any{
		"url": c.URL,
	}
}

type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
}

func (c *NotificationConfig) AsMap() map[string]any {
	return c.ReceiverConfig.AsMap()
}
