package telemetry

import (
	"github.com/newrelic/go-agent/v3/newrelic"
)

// NewRelic contains the New Relic go-agent configuration
type NewRelicConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled" default:"false"`
	AppName string `mapstructure:"appname" yaml:"appname" default:"siren"`
	License string `mapstructure:"license" yaml:"license" default:"____LICENSE_STRING_OF_40_CHARACTERS_____"`
}

func New(c NewRelicConfig) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigEnabled(c.Enabled),
		newrelic.ConfigLicense(c.License),
	)
}
