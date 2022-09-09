package telemetry

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewRelic contains the New Relic go-agent configuration
type NewRelicConfig struct {
	Enabled bool   `yaml:"enabled" mapstructure:"enabled" default:"false"`
	AppName string `yaml:"appname" mapstructure:"appname" default:"siren"`
	License string `yaml:"license" mapstructure:"license"`
}

func New(c NewRelicConfig) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigEnabled(c.Enabled),
		newrelic.ConfigLicense(c.License),
	)
}
