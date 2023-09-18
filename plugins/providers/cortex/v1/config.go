package cortexv1plugin

import (
	_ "embed"
	"fmt"

	"github.com/goto/siren/pkg/httpclient"
	"github.com/goto/siren/plugins/providers/cortex/common"
	"github.com/hashicorp/go-plugin"
)

var (
	//go:embed config/helper.tmpl
	HelperTemplateString string
	//go:embed config/config.goyaml
	ConfigYamlString string
)

type Config struct {
	// https://prometheus.io/docs/alerting/latest/configuration/#route
	GroupWaitDuration      string            `mapstructure:"group_wait" yaml:"group_wait" json:"group_wait" default:"30s"`
	GroupIntervalDuration  string            `mapstructure:"group_interval" yaml:"group_interval" json:"group_interval" default:"5m"`
	RepeatIntervalDuration string            `mapstructure:"repeat_interval" yaml:"repeat_interval" json:"repeat_interval" default:"4h"`
	WebhookBaseAPI         string            `mapstructure:"webhook_base_api" yaml:"webhook_base_api" json:"webhook_base_api" default:"http://localhost:8080/v1beta1/alerts/cortex" validate:"required"`
	HTTPClient             httpclient.Config `mapstructure:"http_client" json:"http_client" yaml:"http_client"`
}

type TemplateConfig struct {
	GroupWaitDuration      string
	GroupIntervalDuration  string `mapstructure:"group_interval" yaml:"group_interval" json:"group_interval" default:"5m"`
	RepeatIntervalDuration string `mapstructure:"repeat_interval" yaml:"repeat_interval" json:"repeat_interval" default:"4h"`
	WebhookURL             string
}

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   common.PluginName,
	MagicCookieValue: fmt.Sprintf("%sv1", common.PluginName),
}
