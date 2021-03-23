package domain

// DBConfig contains the database configuration
type DBConfig struct {
	Host     string `mapstructure:"host" default:"localhost"`
	User     string `mapstructure:"user" default:"postgres"`
	Password string `mapstructure:"password" default:""`
	Name     string `mapstructure:"name" default:"postgres"`
	Port     string `mapstructure:"port" default:"5432"`
	SslMode  string `mapstructure:"sslmode" default:"disable"`
}

// Cortex contains the radar cortex configuration
type CortexConfig struct {
	Host string `mapstructure:"host" default:"http://localhost:8080"`
}

type AlertmanagerConfig struct {
	Address string `mapstructure:"address" default:"localhost:8080"`
}

type NewRelicConfig struct {
	Enabled bool   `mapstructure:"enabled" default:"false"`
	AppName string `mapstructure:"appname" default:"siren"`
	License string `mapstructure:"license"`
}

// Config contains the application configuration
type Config struct {
	Port         int                `mapstructure:"port" default:"8080"`
	DB           DBConfig           `mapstructure:"db"`
	Cortex       CortexConfig       `mapstructure:"cortex"`
	Alertmanager AlertmanagerConfig `mapstructure:"alertmanager"`
	NewRelic     NewRelicConfig     `mapstructure:"newrelic"`
}
