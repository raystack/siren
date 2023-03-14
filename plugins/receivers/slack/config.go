package slack

import (
	"fmt"

	"github.com/goto/siren/pkg/httpclient"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/pkg/secret"
)

// AppConfig is a config loaded when siren is started
type AppConfig struct {
	APIHost    string            `mapstructure:"api_host"`
	Retry      retry.Config      `mapstructure:"retry"`
	HTTPClient httpclient.Config `mapstructure:"http_client"`
}

// SlackCredentialConfig is config that needs to be passed when a new slack
// receiver is being added
type SlackCredentialConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	AuthCode     string `mapstructure:"auth_code"`
}

func (c *SlackCredentialConfig) Validate() error {
	if c.ClientID != "" && c.ClientSecret != "" && c.AuthCode != "" {
		return nil
	}
	return fmt.Errorf("invalid slack credentials, client_id: %s, client_secret: <secret>, auth_code: <secret>", c.ClientID)
}

// ReceiverConfig is a stored config for a slack receiver
type ReceiverConfig struct {
	Token     secret.MaskableString `json:"token" mapstructure:"token"`
	Workspace string                `json:"workspace" mapstructure:"workspace"`
}

func (c *ReceiverConfig) Validate() error {
	if c.Token != "" && c.Workspace != "" {
		return nil
	}
	return fmt.Errorf("invalid slack receiver config, workspace: %s, token: %s", c.Workspace, c.Token)
}

func (c *ReceiverConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"workspace": c.Workspace,
		"token":     c.Token,
	}
}

// ReceiverData is a stored data for a slack receiver
type ReceiverData struct {
	Channels string `json:"channels" mapstructure:"channels"`
}

func (c *ReceiverData) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"channels": c.Channels,
	}
}

// SubscriptionConfig is a stored config for a subscription of a slack receiver
type SubscriptionConfig struct {
	ChannelName string `json:"channel_name" mapstructure:"channel_name"`
	ChannelType string `json:"channel_type" mapstructure:"channel_type"`
}

func (c *SubscriptionConfig) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"channel_name": c.ChannelName,
		"channel_type": c.ChannelType,
	}
}

// NotificationConfig has all configs needed to send notification
type NotificationConfig struct {
	ReceiverConfig     `mapstructure:",squash"`
	SubscriptionConfig `mapstructure:",squash"`
}

// Validate validates whether notification config contains required fields or not
// channel_name is not mandatory because in NotifyToReceiver flow, channel_name
// is being passed from the request (not from the config)
func (c *NotificationConfig) Validate() error {
	if c.Token != "" && c.Workspace != "" {
		return nil
	}
	return fmt.Errorf("invalid slack notification config, workspace: %s, token: %s, channel_name: %s", c.Workspace, c.Token, c.ChannelName)
}

func (c *NotificationConfig) AsMap() map[string]interface{} {
	notificationMap := make(map[string]interface{})

	for k, v := range c.ReceiverConfig.AsMap() {
		notificationMap[k] = v
	}

	for k, v := range c.SubscriptionConfig.AsMap() {
		notificationMap[k] = v
	}

	return notificationMap
}
