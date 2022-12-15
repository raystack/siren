package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/idempotency"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/pgc"
)

const idempotencyInsertQuery = `
INSERT INTO idempotencies (scope, key, success, created_at, updated_at)
    VALUES ($1, $2, false, now(), now()) ON CONFLICT (scope, key) DO UPDATE SET scope=$1, updated_at=now()
RETURNING *
`

const idempotencyUpdateQuery = `
UPDATE idempotencies SET success=$2, updated_at=now()
  WHERE id=$1
`

const idempotencyDeleteTemplateQuery = `
DELETE FROM idempotencies WHERE now() - interval '%d seconds' > updated_at
`

// IdempotencyRepository talks to the store to read or insert idempotency keys
type IdempotencyRepository struct {
	client    *pgc.Client
	tableName string
}

// NewIdempotencyRepository returns repository struct
func NewIdempotencyRepository(client *pgc.Client) *IdempotencyRepository {
	return &IdempotencyRepository{client, "idempotencies"}
}

func (r *IdempotencyRepository) InsertOnConflictReturning(ctx context.Context, scope, key string) (*idempotency.Idempotency, error) {
	var idempotencyModel model.Idempotency
	if err := r.client.QueryRowxContext(ctx, pgc.OpInsert, r.tableName, idempotencyInsertQuery,
		scope, key,
	).StructScan(&idempotencyModel); err != nil {
		return nil, pgc.CheckError(err)
	}

	return idempotencyModel.ToDomain(), nil
}

func (r *IdempotencyRepository) UpdateSuccess(ctx context.Context, id uint64, success bool) error {
	rows, err := r.client.ExecContext(ctx, pgc.OpUpdate, r.tableName, idempotencyUpdateQuery,
		id, success,
	)
	if err != nil {
		return err
	}

	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *IdempotencyRepository) Delete(ctx context.Context, filter idempotency.Filter) error {

	if filter.TTL == 0 {
		return errors.ErrInvalid.WithCausef("cannot delete with ttl 0")
	}

	ttlInSecond := int(filter.TTL.Seconds())

	rows, err := r.client.ExecContext(ctx, pgc.OpDelete, r.tableName, fmt.Sprintf(idempotencyDeleteTemplateQuery, ttlInSecond))
	if err != nil {
		return err
	}

	ra, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if ra == 0 {
		return errors.ErrNotFound
	}

	return nil
}
