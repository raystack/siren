package postgres

import (
	"fmt"

	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New returns the database instance
func New(c Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=%s password=%s ",
		c.Host,
		c.User,
		c.Name,
		c.Port,
		c.SSLMode,
		c.Password,
	)

	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(getLogLevelFromString(c.LogLevel))})
	if err != nil {
		panic(err)
	}

	return db, err
}

func getLogLevelFromString(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
