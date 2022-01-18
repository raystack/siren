package config

import (
	"errors"
	"fmt"

	"github.com/odpf/salt/config"
	"github.com/odpf/siren/domain"
)

// LoadConfig returns application configuration
func LoadConfig(configFile string) (*domain.Config, error) {
	var cfg domain.Config
	loader := config.NewLoader(config.WithFile(configFile))

	if err := loader.Load(&cfg); err != nil {
		if errors.As(err, &config.ConfigFileNotFoundError{}) {
			fmt.Println(err)
			return &cfg, nil
		}
		return nil, err
	}
	return &cfg, nil
}
