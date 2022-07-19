package postgres

// Config contains the database configuration
type Config struct {
	Host     string `mapstructure:"host" default:"localhost"`
	User     string `mapstructure:"user" default:"postgres"`
	Password string `mapstructure:"password" default:""`
	Name     string `mapstructure:"name" default:"postgres"`
	Port     string `mapstructure:"port" default:"5432"`
	SSLMode  string `mapstructure:"sslmode" default:"disable"`
}
