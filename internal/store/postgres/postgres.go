package postgres

import (
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	errDuplicateKey        = errors.New("duplicate key")
	errCheckViolation      = errors.New("check constraint violation")
	errForeignKeyViolation = errors.New("foreign key violation")
)

type Client struct {
	db     *gorm.DB
	logger log.Logger
}

// New returns the database instance// NewClient initializes database connection
func NewClient(logger log.Logger, c Config) (*Client, error) {
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
		return nil, err
	}

	return &Client{
		db:     db,
		logger: logger,
	}, nil
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

func checkPostgresError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("%w [%s]", errDuplicateKey, pgErr.Detail)
		case pgerrcode.CheckViolation:
			return fmt.Errorf("%w [%s]", errCheckViolation, pgErr.Detail)
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf("%w [%s]", errForeignKeyViolation, pgErr.Detail)
		}
	}
	return err
}

func (c *Client) Migrate() error {
	c.logger.Info("migrating postgres...")
	err := c.db.AutoMigrate(
		&model.Alert{},
		&model.Namespace{},
		&model.Provider{},
		&model.Receiver{},
		&model.Rule{},
		&model.Subscription{},
		&model.Template{})
	if err != nil {
		return err
	}
	c.logger.Info("migration done.")
	return nil
}

// GetDB is a helper function for purpose only
// don't use it in the main code
func (c *Client) GetDB() *gorm.DB {
	return c.db
}
