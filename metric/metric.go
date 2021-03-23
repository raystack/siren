package metric

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/siren/domain"
)

// RunServer runs the application server
func New(c *domain.NewRelicConfig) (*newrelic.Application, error) {
	return newrelic.NewApplication(
		newrelic.ConfigAppName(c.AppName),
		newrelic.ConfigEnabled(c.Enabled),
		newrelic.ConfigLicense(c.License),
	)
}
