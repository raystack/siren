package providers

import (
	"github.com/raystack/siren/plugins/providers/cortex"
)

type Config struct {
	Cortex cortex.AppConfig `mapstructure:"cortex"`
}
