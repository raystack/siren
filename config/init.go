package config

import (
	"os"

	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v2"
)

func Init(configFile string) error {
	cfg := &Config{}

	defaults.SetDefaults(cfg)

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0655); err != nil {
		return err
	}

	return nil
}
