package cortex

import (
	_ "embed"

	"github.com/goto/siren/pkg/httpclient"
)

var (
	//go:embed config/helper.tmpl
	HelperTemplateString string
	//go:embed config/config.goyaml
	ConfigYamlString string
)

type AppConfig struct {
	// https://prometheus.io/docs/alerting/latest/configuration/#route
	GroupWaitDuration      string            `mapstructure:"group_wait" yaml:"group_wait" default:"30s"`
	GroupIntervalDuration  string            `mapstructure:"group_interval" yaml:"group_interval" default:"5m"`
	RepeatIntervalDuration string            `mapstructure:"repeat_interval" yaml:"repeat_interval" default:"4h"`
	WebhookBaseAPI         string            `mapstructure:"webhook_base_api" yaml:"webhook_base_api" default:"http://localhost:8080/v1beta1/alerts/cortex"`
	HTTPClient             httpclient.Config `mapstructure:"http_client" yaml:"http_client"`
}

type TemplateConfig struct {
	GroupWaitDuration      string
	GroupIntervalDuration  string `mapstructure:"group_interval" yaml:"group_interval" default:"5m"`
	RepeatIntervalDuration string `mapstructure:"repeat_interval" yaml:"repeat_interval" default:"4h"`
	WebhookURL             string
}
