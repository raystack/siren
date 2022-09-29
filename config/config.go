package config

import (
	_ "embed"
	"fmt"

	"github.com/odpf/salt/config"
	"github.com/odpf/salt/db"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
	"github.com/odpf/siren/plugins/providers/cortex"
	"github.com/odpf/siren/plugins/receivers"
)

// Load returns application configuration
func Load(configFile string) (Config, error) {
	var cfg Config
	loader := config.NewLoader(config.WithFile(configFile))

	if err := loader.Load(&cfg); err != nil {
		if errors.As(err, &config.ConfigFileNotFoundError{}) {
			fmt.Println(err)
			return cfg, nil
		}
		return cfg, err
	}

	return cfg, nil
}

type Log struct {
	Level         string `mapstructure:"level" default:"info"`
	GCPCompatible bool   `mapstructure:"gcp_compatible" default:"true"`
}

// Config contains the application configuration
type Config struct {
	DB           db.Config                `mapstructure:"db"`
	Cortex       cortex.Config            `mapstructure:"cortex"`
	NewRelic     telemetry.NewRelicConfig `mapstructure:"newrelic"`
	Service      server.Config            `mapstructure:"service"`
	Log          Log                      `mapstructure:"log"`
	Receivers    receivers.Config         `mapstructure:"receivers"`
	Notification notification.Config      `mapstructure:"notification"`
}
