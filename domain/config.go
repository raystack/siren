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

// CortexConfig contains the cortex configuration
type CortexConfig struct {
	Address string `mapstructure:"address" default:"http://localhost:8080"`
}

// NewRelic contains the New Relic go-agent configuration
type NewRelicConfig struct {
	Enabled bool   `mapstructure:"enabled" default:"false"`
	AppName string `mapstructure:"appname" default:"siren"`
	License string `mapstructure:"license"`
}

// Config contains the application configuration
type Config struct {
	Port     int            `mapstructure:"port" default:"8080"`
	DB       DBConfig       `mapstructure:"db"`
	Cortex   CortexConfig   `mapstructure:"cortex"`
	NewRelic NewRelicConfig `mapstructure:"newrelic"`
}
