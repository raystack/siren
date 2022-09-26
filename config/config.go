package config

import (
	_ "embed"
	"fmt"

	"github.com/odpf/salt/config"
	"github.com/odpf/salt/db"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
	"github.com/odpf/siren/plugins/providers/cortex"
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
	Level         string `yaml:"level" mapstructure:"level" default:"info"`
	GCPCompatible bool   `yaml:"gcp_compatible" mapstructure:"gcp_compatible" default:"true"`
}

type SlackApp struct {
	ClientID     string `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
}

// Config contains the application configuration
type Config struct {
	DB            db.Config                `mapstructure:"db"`
	Cortex        cortex.Config            `mapstructure:"cortex"`
	NewRelic      telemetry.NewRelicConfig `mapstructure:"newrelic"`
	Service       server.Config            `mapstructure:"service"`
	Log           Log                      `mapstructure:"log"`
	SlackApp      SlackApp                 `mapstructure:"slack_app"`
	EncryptionKey string                   `mapstructure:"encryption_key"`
}
