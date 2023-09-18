package cortexv1plugin

import "errors"

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

func (a Alert) Validate() error {
	if _, ok := a.Labels["severity"]; !ok {
		return errors.New("'severity' label is missing")
	}

	if _, ok := a.Annotations["resource"]; !ok {
		return errors.New("'resource' annotation is missing")
	}

	if _, ok := a.Annotations["template"]; !ok {
		return errors.New("'template' annotation is missing")
	}

	if _, ok := a.Annotations["metric_value"]; !ok {
		return errors.New("'metric_value' annotation is missing")
	}

	if _, ok := a.Annotations["metric_name"]; !ok {
		return errors.New("'metric_name' annotation is missing")
	}

	return nil
}
