package config

import (
	_ "embed"
	"fmt"

	"github.com/odpf/salt/config"
	"github.com/odpf/salt/db"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/telemetry"
)

var (
	//go:embed prometheus_alert_manager_helper.tmpl
	promAMHelperTemplateString string
	//go:embed prometheus_alert_manager_config.goyaml
	promAMConfigYamlString string
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
	cfg.Cortex.PrometheusAlertManagerConfigYaml = promAMConfigYamlString
	cfg.Cortex.PrometheusAlertManagerHelperTemplate = promAMHelperTemplateString
	return cfg, nil
}

type LogConfig struct {
	Level string `yaml:"level" mapstructure:"level" default:"info"`
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
	Log           LogConfig                `mapstructure:"log"`
	SlackApp      SlackApp                 `mapstructure:"slack_app"`
	EncryptionKey string                   `mapstructure:"encryption_key"`
}
