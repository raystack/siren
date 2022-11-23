package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/store/postgres/migrations"
	"github.com/odpf/siren/pkg/errors"
)

var (
	transactionContextKey = struct{}{}

	errDuplicateKey        = errors.New("duplicate key")
	errCheckViolation      = errors.New("check constraint violation")
	errForeignKeyViolation = errors.New("foreign key violation")
)

type Client struct {
	db     *db.Client
	logger log.Logger
	// postgresTracer *telemetry.PostgresSpan
}

// NewClient wraps salt/db client
func NewClient(logger log.Logger, dbc *db.Client) (*Client, error) {
	if dbc == nil {
		return nil, errors.New("error creating postgres client: nil db client")
	}

	// postgresTracer, err := telemetry.InitPostgresSpan(
	// 	"public",
	// 	dbc.ConnectionURL(),
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &Client{
		db:     dbc,
		logger: logger,
		// postgresTracer: postgresTracer,
	}, nil
}

func checkPostgresError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return fmt.Errorf("%w [%s]", errDuplicateKey, pqErr.Detail)
		case "check_violation":
			return fmt.Errorf("%w [%s]", errCheckViolation, pqErr.Detail)
		case "foreign_key_violation":
			return fmt.Errorf("%w [%s]", errForeignKeyViolation, pqErr.Detail)
		}
	}
	return err
}

func Migrate(cfg db.Config) error {
	if err := db.RunMigrations(cfg, migrations.FS, migrations.ResourcePath); err != nil {
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
