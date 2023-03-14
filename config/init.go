package config

import (
	"os"

	"github.com/goto/siren/core/receiver"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

func Init(configFile string) error {
	cfg := &Config{}

	defaults.SetDefaults(cfg)

	cfg.DB.Driver = "postgres"
	cfg.DB.URL = "postgres://postgres:@localhost:5432/siren_development?sslmode=disable"

	if len(cfg.Notification.MessageHandler.ReceiverTypes) == 0 {
		cfg.Notification.MessageHandler.ReceiverTypes = receiver.SupportedTypes
	}

	if len(cfg.Notification.DLQHandler.ReceiverTypes) == 0 {
		cfg.Notification.DLQHandler.ReceiverTypes = receiver.SupportedTypes
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}
