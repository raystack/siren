package cortex

import _ "embed"

var (
	//go:embed config/helper.tmpl
	HelperTemplateString string
	//go:embed config/config.goyaml
	ConfigYamlString string
)

type AppConfig struct {
	GroupWaitDuration string `mapstructure:"group_wait" yaml:"group_wait" default:"30s"`
	WebhookBaseAPI    string `mapstructure:"webhook_base_api" yaml:"webhook_base_api" default:"http://localhost:8080/v1beta1/alerts/cortex"`
}

type TemplateConfig struct {
	GroupWaitDuration string
	WebhookURL        string
}
