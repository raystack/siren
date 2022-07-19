package postgres

import (
	"context"
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
	transactionContextKey = struct{}{}

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
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	case "debug":
		return gormlogger.Info
	default:
		return gormlogger.Silent
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

func (c *Client) WithTransaction(ctx context.Context) context.Context {
	tx := c.db.Begin()
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (c *Client) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Rollback()
		if tx.Error != nil {
			return c.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Commit()
		if tx.Error != nil {
			return c.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) GetDB(ctx context.Context) *gorm.DB {
	db := c.db
	if tx := extractTransaction(ctx); tx != nil {
		db = tx
	}
	return db
}

func extractTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); !ok {
		return nil
	} else {
		return tx
	}
}
