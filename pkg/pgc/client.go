package pgc

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
	"github.com/odpf/siren/pkg/telemetry"
	"go.opencensus.io/trace"
)

const (
	OpInsert    = "INSERT"
	OpSelectAll = "SELECT_ALL"
	OpSelect    = "SELECT"
	OpUpdate    = "UPDATE"
	OpDelete    = "DELETE"
)

var (
	transactionContextKey = struct{}{}

	ErrDuplicateKey        = errors.New("duplicate key")
	ErrCheckViolation      = errors.New("check constraint violation")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

type Client struct {
	db             *db.Client
	logger         log.Logger
	postgresTracer *telemetry.PostgresTracer
}

// NewClient wraps salt/db client
func NewClient(logger log.Logger, dbc *db.Client) (*Client, error) {
	if dbc == nil {
		return nil, errors.New("error creating postgres client: nil db client")
	}

	postgresTracer, err := telemetry.NewPostgresTracer(
		dbc.ConnectionURL(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:             dbc,
		logger:         logger,
		postgresTracer: postgresTracer,
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

func Migrate(cfg db.Config) error {
	if err := db.RunMigrations(cfg, migrations.FS, migrations.ResourcePath); err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryRowxContext(ctx context.Context, op string, tableName string, query string, args ...interface{}) *sqlx.Row {
	ctx, span := c.postgresTracer.StartSpan(ctx, op, tableName,
		trace.StringAttribute("db.statement", query),
	)
	defer span.End()
	sqlxRow := c.GetDB(ctx).QueryRowxContext(ctx, query, args...)
	if sqlxRow.Err() != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: sqlxRow.Err().Error(),
		})
	}
	return sqlxRow
}

func (c *Client) QueryxContext(ctx context.Context, op string, tableName string, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx, span := c.postgresTracer.StartSpan(ctx, op, tableName,
		trace.StringAttribute("db.statement", query),
	)
	defer span.End()
	sqlxRow, err := c.GetDB(ctx).QueryxContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}
	return sqlxRow, err
}

func (c *Client) GetContext(ctx context.Context, op string, tableName string, dest interface{}, query string, args ...interface{}) error {
	ctx, span := c.postgresTracer.StartSpan(ctx, op, tableName,
		trace.StringAttribute("db.statement", query),
	)
	defer span.End()

	if err := c.GetDB(ctx).QueryRowxContext(ctx, query, args...).StructScan(dest); err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
		return err
	}

	return nil
}

func (c *Client) ExecContext(ctx context.Context, op string, tableName string, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := c.postgresTracer.StartSpan(ctx, op, tableName,
		trace.StringAttribute("db.statement", query),
	)
	defer span.End()

	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
		return nil, err
	}

	return res, nil
}

func (c *Client) NamedExecContext(ctx context.Context, op string, tableName string, query string, arg interface{}) (sql.Result, error) {
	ctx, span := c.postgresTracer.StartSpan(ctx, op, tableName,
		trace.StringAttribute("db.statement", query),
	)
	defer span.End()

	res, err := c.db.NamedExecContext(ctx, query, arg)
	if err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
		return nil, err
	}

	return res, nil
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
