package postgres

import (
	"fmt"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/store/model"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New returns the database instance
func New(logger log.Logger, c Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=%s password=%s ",
		c.Host,
		c.User,
		c.Name,
		c.Port,
		c.SSLMode,
		c.Password,
	)

	db, err := gorm.Open(gormpg.Open(dsn), &gorm.Config{Logger: gormlogger.Default.LogMode(getLogLevelFromString(logger.Level()))})
	if err != nil {
		panic(err)
	}

	return db, err
}

func getLogLevelFromString(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Info
	}
}

func Migrate(logger log.Logger, db *gorm.DB) error {
	logger.Info("migrating alert table...")
	err := db.AutoMigrate(&model.Alert{})
	if err != nil {
		return err
	}
	logger.Info("migrating namespace table...")
	err = db.AutoMigrate(&model.Namespace{})
	if err != nil {
		return err
	}
	logger.Info("migrating provider table...")
	err = db.AutoMigrate(&model.Provider{})
	if err != nil {
		return err
	}
	logger.Info("migrating receiver table...")
	err = db.AutoMigrate(&model.Receiver{})
	if err != nil {
		return err
	}
	logger.Info("migrating rule table...")
	err = db.AutoMigrate(&model.Rule{})
	if err != nil {
		return err
	}
	logger.Info("migrating subscription table...")
	err = db.AutoMigrate(&model.Subscription{})
	if err != nil {
		return err
	}
	logger.Info("migrating template table...")
	err = db.AutoMigrate(&model.Template{})
	if err != nil {
		return err
	}
	logger.Info("migration done.")
	return nil
}
