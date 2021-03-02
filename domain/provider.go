package domain

import "time"

// DBConfig contains the database configuration
type DBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name" default:"postgres"`
	Port     string `mapstructure:"port" default:"5432"`
	SslMode  string `mapstructure:"sslmode"`
}

// Config contains the application configuration
type Config struct {
	Port int      `mapstructure:"port" default:"8080"`
	DB   DBConfig `mapstructure:"db"`
}

type Template struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Tags      string    `gorm:"type:text[]" json:"tags"`
}
