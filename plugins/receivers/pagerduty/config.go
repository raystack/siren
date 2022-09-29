package pagerduty

import (
	"fmt"
)

type ReceiverConfig struct {
	ServiceKey string `mapstructure:"service_key"`
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
	ReceiverConfig
}

func (c *NotificationConfig) AsMap() map[string]interface{} {
	return c.ReceiverConfig.AsMap()
}
