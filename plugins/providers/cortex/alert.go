package cortex

import (
	"github.com/odpf/siren/core/alert"
)

// GroupAlert contract is cortex/prometheus webhook_config contract
// https://prometheus.io/docs/alerting/latest/configuration/#webhook_config
type GroupAlert struct {
	GroupKey    string  `mapstructure:"groupKey"`
	ExternalURL string  `mapstructure:"externalUrl"`
	Version     string  `mapstructure:"version"`
	Alerts      []Alert `mapstructure:"alerts"`
}

type Alert struct {
	Status       string            `mapstructure:"status"`
	Annotations  map[string]string `mapstructure:"annotations"`
	Labels       map[string]string `mapstructure:"labels"`
	GeneratorURL string            `mapstructure:"generatorURL"`
	Fingerprint  string            `mapstructure:"fingerprint"`
	StartsAt     string            `mapstructure:"startsAt"`
	EndsAt       string            `mapstructure:"endsAt"`
}

func isValidCortexAlert(alrt alert.Alert) bool {
	return !(alrt.ResourceName == "" || alrt.Rule == "" ||
		alrt.MetricValue == "" || alrt.MetricName == "" ||
		alrt.Severity == "")
}
