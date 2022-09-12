package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/errors"
)

//go:embed migrations/*.sql
var fs embed.FS

var (
	transactionContextKey = struct{}{}

	errDuplicateKey        = errors.New("duplicate key")
	errCheckViolation      = errors.New("check constraint violation")
	errForeignKeyViolation = errors.New("foreign key violation")
)

type Client struct {
	db     *db.Client
	logger log.Logger
}

// NewClient wraps salt/db client
func NewClient(logger log.Logger, dbc *db.Client) (*Client, error) {
	if dbc == nil {
		return nil, errors.New("error creating postgres client: nil db client")
	}

	return &Client{
		db:     dbc,
		logger: logger,
	}, nil
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

func Migrate(cfg db.Config) error {
	if err := db.RunMigrations(cfg, fs, "migrations"); err != nil {
		return err
	}
	return nil
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

func extractTransaction(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(transactionContextKey).(*sqlx.Tx); !ok {
		return nil
	} else {
		return tx
	}
}
