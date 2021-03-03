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

// Config contains the application configuration
type Config struct {
	Port int      `mapstructure:"port" default:"8080"`
	DB   DBConfig `mapstructure:"db"`
}
