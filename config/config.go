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
	"github.com/odpf/siren/plugins/providers"
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
	Level         string `mapstructure:"level" yaml:"level" default:"info"`
	GCPCompatible bool   `mapstructure:"gcp_compatible" yaml:"gcp_compatible" default:"true"`
}

// Config contains the application configuration
type Config struct {
	DB           db.Config           `mapstructure:"db"`
	Telemetry    telemetry.Config    `mapstructure:"telemetry" yaml:"telemetry"`
	Service      server.Config       `mapstructure:"service" yaml:"service"`
	Log          Log                 `mapstructure:"log" yaml:"log"`
	Providers    providers.Config    `mapstructure:"providers" yaml:"providers"`
	Receivers    receivers.Config    `mapstructure:"receivers" yaml:"receivers"`
	Notification notification.Config `mapstructure:"notification" yaml:"notification"`
}
