package pagerduty

import "github.com/goto/siren/pkg/secret"

// https://developer.pagerduty.com/docs/ZG9jOjExMDI5NTc4-send-a-v1-event
type MessageV1 struct {
	ServiceKey  secret.MaskableString `mapstructure:"service_key" yaml:"service_key,omitempty" json:"service_key,omitempty"`
	EventType   string                `mapstructure:"event_type" yaml:"event_type,omitempty" json:"event_type,omitempty"`
	IncidentKey string                `mapstructure:"incident_key" yaml:"incident_key,omitempty" json:"incident_key,omitempty"`
	Description string                `mapstructure:"description" yaml:"description,omitempty" json:"description,omitempty"`
	Details     map[string]any        `mapstructure:"details" yaml:"details,omitempty" json:"details,omitempty"`
	Client      string                `mapstructure:"client" yaml:"client,omitempty" json:"client,omitempty"`
	ClientURL   string                `mapstructure:"client_url" yaml:"client_url,omitempty" json:"client_url,omitempty"`
	Contexts    []Context             `mapstructure:"contexts" yaml:"contexts,omitempty" json:"contexts,omitempty"`
}

type Context struct {
	Type string `mapstructure:"type" yaml:"type,omitempty" json:"type"`
	Src  string `mapstructure:"src" yaml:"src,omitempty" json:"src"`
	Href string `mapstructure:"href" yaml:"href,omitempty" json:"href"`
	Text string `mapstructure:"text" yaml:"text,omitempty" json:"text"`
	Alt  string `mapstructure:"alt" yaml:"alt,omitempty" json:"alt"`
}
