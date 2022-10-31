package providers

import (
	"github.com/odpf/siren/plugins/providers/cortex"
)

type Config struct {
	Cortex cortex.AppConfig `mapstructure:"cortex"`
}
