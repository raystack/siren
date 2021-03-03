package store

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"gorm.io/gorm/logger"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New returns the database instance
func New(c *domain.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=%s password=%s ",
		c.Host,
		c.User,
		c.Name,
		c.Port,
		c.SslMode,
		c.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Panic(err)
	}

	return db, err
}

// Migrate auto migrate models
func Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
