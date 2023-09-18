package config

import (
	_ "embed"
	"fmt"

	"github.com/goto/salt/config"
	"github.com/goto/salt/db"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/server"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/telemetry"
	"github.com/goto/siren/plugins"
	"github.com/goto/siren/plugins/receivers"
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
	Providers    plugins.Config      `mapstructure:"providers" yaml:"providers"`
	Receivers    receivers.Config    `mapstructure:"receivers" yaml:"receivers"`
	Notification notification.Config `mapstructure:"notification" yaml:"notification"`
}
