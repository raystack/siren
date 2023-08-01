package slackchannel

import (
	"fmt"

	"github.com/goto/siren/plugins/receivers/slack"
)

// ReceiverConfig is a stored config for a slack receiver
type ReceiverConfig struct {
	SlackReceiverConfig slack.ReceiverConfig `mapstructure:",squash"`
	ChannelName         string               `json:"channel_name" mapstructure:"channel_name"`
	ChannelType         string               `json:"channel_type" mapstructure:"channel_type"`
}

func (c *ReceiverConfig) Validate() error {
	if c.ChannelName == "" {
		return fmt.Errorf("invalid slack_channel receiver config, channel_name can't be empty")
	}
	return nil
}

func (c *ReceiverConfig) AsMap() map[string]any {
	return map[string]any{
		"token":        c.SlackReceiverConfig.Token,
		"workspace":    c.SlackReceiverConfig.Workspace,
		"channel_name": c.ChannelName,
		"channel_type": c.ChannelType,
	}
}

// NotificationConfig has all configs needed to send notification
type NotificationConfig struct {
	ReceiverConfig `mapstructure:",squash"`
	// SubscriptionConfig `mapstructure:",squash"`
}

// Validate validates whether notification config contains required fields or not
// channel_name is not mandatory because in NotifyToReceiver flow, channel_name
// is being passed from the request (not from the config)
func (c *NotificationConfig) Validate() error {
	if err := c.ReceiverConfig.SlackReceiverConfig.Validate(); err != nil {
		return err
	}
	if err := c.ReceiverConfig.Validate(); err != nil {
		return err
	}
	return nil
}

func (c *NotificationConfig) AsMap() map[string]any {
	notificationMap := make(map[string]any)

	for k, v := range c.ReceiverConfig.AsMap() {
		notificationMap[k] = v
	}

	return notificationMap
}
