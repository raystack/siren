package receivers

import (
	"github.com/odpf/siren/plugins/receivers/httpreceiver"
	"github.com/odpf/siren/plugins/receivers/pagerduty"
	"github.com/odpf/siren/plugins/receivers/slack"
)

type Config struct {
	Slack        slack.AppConfig        `mapstructure:"slack" yaml:"slack"`
	Pagerduty    pagerduty.AppConfig    `mapstructure:"pagerduty" yaml:"pagerduty"`
	HTTPReceiver httpreceiver.AppConfig `mapstructure:"http" yaml:"http"`
}
