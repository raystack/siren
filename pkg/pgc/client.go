package pgc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goto/salt/db"
	"github.com/goto/salt/log"
	"github.com/goto/siren/pkg/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	transactionContextKey = struct{}{}

	ErrDuplicateKey        = errors.New("duplicate key")
	ErrCheckViolation      = errors.New("check constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type Client struct {
	db     *db.Client
	logger log.Logger
}

// NewClient wraps salt/db client
func NewClient(logger log.Logger, cfg db.Config) (*Client, error) {
	driverName, err := otelsql.Register(
		cfg.Driver,
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, fmt.Errorf("new pgq processor: %w", err)
	}

	// backup origin driver name that is going to be overrided
	originDriverName := cfg.Driver
	cfg.Driver = driverName

	dbClient, err := db.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres client: %w", err)
	}

	sqlDB := dbClient.DB.DB

	// need to use NewDb if we want to use NamedContext with otelsql
	// ref: https://github.com/nhatthm/otelsql?tab=readme-ov-file#jmoironsqlx
	wrappedSQLxDB := sqlx.NewDb(sqlDB, originDriverName)

	if err := otelsql.RecordStats(
		sqlDB,
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
		otelsql.WithInstanceName(dbClient.Host()),
	); err != nil {
		return nil, err
	}

	dbClient.DB = wrappedSQLxDB

	return &Client{
		db:     dbClient,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func CheckError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return fmt.Errorf("%w [%s]", ErrDuplicateKey, pqErr.Detail)
		case "check_violation":
			return fmt.Errorf("%w [%s]", ErrCheckViolation, pqErr.Detail)
		case "foreign_key_violation":
			return fmt.Errorf("%w [%s]", ErrForeignKeyViolation, pqErr.Detail)
		}
	}
	return err
}

func (c *Client) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return c.GetDB(ctx).QueryRowxContext(ctx, query, args...)
}

func (c *Client) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.GetDB(ctx).QueryxContext(ctx, query, args...)
}

func (c *Client) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.GetDB(ctx).QueryRowxContext(ctx, query, args...).StructScan(dest)
}

func (c *Client) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

func (c *Client) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return c.db.NamedExecContext(ctx, query, arg)
}

func (c *Client) WithTransaction(ctx context.Context, opts *sql.TxOptions) context.Context {
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (c *Client) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	return errors.New("no transaction")
}

func (c *Client) GetDB(ctx context.Context) sqlx.QueryerContext {
	if tx := extractTransaction(ctx); tx != nil {
		return tx
	}
	return c.db
}

func (c *Client) GetSaltDB(ctx context.Context) *db.Client {
	return c.db
}

func extractTransaction(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(transactionContextKey).(*sqlx.Tx); !ok {
		return nil
	} else {
		return tx
	}
}
