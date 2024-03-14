package postgres

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
)

var idempotenctListQueryBuilder = sq.Select(
	"scope",
	"key",
	"notification_id",
	"created_at",
	"updated_at",
).From("idempotencies")

const idempotencyInsertQuery = `
INSERT INTO idempotencies (scope, key, notification_id, created_at, updated_at)
    VALUES ($1, $2, $3, now(), now()) RETURNING *
`

const idempotencyDeleteTemplateQuery = `
DELETE FROM idempotencies WHERE now() - interval '%d seconds' > updated_at
`

// IdempotencyRepository talks to the store to read or insert idempotency keys
type IdempotencyRepository struct {
	client *pgc.Client
}

// NewIdempotencyRepository returns repository struct
func NewIdempotencyRepository(client *pgc.Client) *IdempotencyRepository {
	return &IdempotencyRepository{client}
}

func (r *IdempotencyRepository) Create(ctx context.Context, scope, key, notificationID string) (*notification.Idempotency, error) {
	if scope == "" || key == "" {
		return nil, errors.ErrInvalid.WithMsgf("scope or key cannot be empty")
	}
	var idempotencyModel model.Idempotency
	if err := r.client.QueryRowxContext(ctx, idempotencyInsertQuery,
		scope, key, notificationID,
	).StructScan(&idempotencyModel); err != nil {
		return nil, pgc.CheckError(err)
	}

	return idempotencyModel.ToDomain(), nil
}

func (r *IdempotencyRepository) Check(ctx context.Context, scope, key string) (*notification.Idempotency, error) {
	var idempotencyModel model.Idempotency

	query, args, err := idempotenctListQueryBuilder.
		Where("scope = ?", scope).
		Where("key = ?", key).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.client.GetContext(ctx, &idempotencyModel, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	return idempotencyModel.ToDomain(), nil
}

func (r *IdempotencyRepository) Delete(ctx context.Context, filter notification.IdempotencyFilter) error {

	if filter.TTL == 0 {
		return errors.ErrInvalid.WithCausef("cannot delete with ttl 0")
	}

	ttlInSecond := int(filter.TTL.Seconds())

	rows, err := r.client.ExecContext(ctx, fmt.Sprintf(idempotencyDeleteTemplateQuery, ttlInSecond))
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
