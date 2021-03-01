package app

import (
	"fmt"
	"github.com/jeremywohl/flatten"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

// DBConfig contains the database configuration
type DBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name" default:"postgres"`
	Port     string `mapstructure:"port"`
	SslMode  string `mapstructure:"sslmode"`
}

// Config contains the application configuration
type Config struct {
	Port int      `mapstructure:"port"`
	DB   DBConfig `mapstructure:"db"`
}

// LoadConfig returns application configuration
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("viper read config error %v", err)
	}

	err, configKeys := getFlattenedStructKeys(Config{})
	if err != nil {
		fmt.Errorf("Unable to get config keys : %s\n", err)
	}

	// Bind each conf fields to environment vars
	for key := range configKeys {
		err := viper.BindEnv(configKeys[key])
		if err != nil {
			fmt.Errorf("Unable to bind env var: %s\n", err)
		}
	}

	var config Config
	defaults.SetDefaults(&config)

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Errorf("unable to unmarshal config to struct: %v\n", err)
	}
	return &config
}

func getFlattenedStructKeys(config Config) (error, []string) {
	var structMap map[string]interface{}
	err := mapstructure.Decode(config, &structMap)
	if err != nil {
		return err, nil
	}

	flat, err := flatten.Flatten(structMap, "", flatten.DotStyle)
	if err != nil {
		return err, nil
	}

	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}

	return nil, keys
}
