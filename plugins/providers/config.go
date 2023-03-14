package providers

import (
	"github.com/goto/siren/plugins/providers/cortex"
)

type Config struct {
	Cortex cortex.AppConfig `mapstructure:"cortex"`
}
