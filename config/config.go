package config

import (
	_ "embed"
	"fmt"

	"github.com/odpf/salt/config"
	"github.com/odpf/siren/internal/server"
	"github.com/odpf/siren/internal/store/postgres"
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

// LoadConfig returns application configuration
func LoadConfig(configFile string) (Config, error) {
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
	Level string `mapstructure:"level" default:"info"`
}

type SlackApp struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// Config contains the application configuration
type Config struct {
	DB            postgres.Config          `mapstructure:"db"`
	Cortex        cortex.Config            `mapstructure:"cortex"`
	NewRelic      telemetry.NewRelicConfig `mapstructure:"newrelic"`
	SirenService  server.Config            `mapstructure:"siren_service"`
	Log           LogConfig                `mapstructure:"log"`
	SlackApp      SlackApp                 `mapstructure:"slack_app"`
	EncryptionKey string                   `mapstructure:"encryption_key"`
}
