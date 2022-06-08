package telemetry

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/siren/domain"
)

func New(c *domain.NewRelicConfig) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigEnabled(c.Enabled),
		newrelic.ConfigLicense(c.License),
	)
}
