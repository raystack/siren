package postgres

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const subscriptionInsertQuery = `
INSERT INTO subscriptions (namespace_id, urn, receiver, match, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
RETURNING *
`

const subscriptionUpdateQuery = `
UPDATE subscriptions SET namespace_id=$2, urn=$3, receiver=$4, match=$5, updated_at=now()
WHERE id = $1
RETURNING *
`

const subscriptionDeleteQuery = `
DELETE from subscriptions where id=$1
`

var subscriptionListQueryBuilder = sq.Select(
	"id",
	"namespace_id",
	"urn",
	"receiver",
	"match",
	"created_at",
	"updated_at",
).From("subscriptions")

// SubscriptionRepository talks to the store to read or insert data
type SubscriptionRepository struct {
	client *Client
}

// NewSubscriptionRepository returns SubscriptionRepository struct
func NewSubscriptionRepository(client *Client) *SubscriptionRepository {
	return &SubscriptionRepository{
		client: client,
	}
}

func (r *SubscriptionRepository) List(ctx context.Context, flt subscription.Filter) ([]subscription.Subscription, error) {
	var queryBuilder = subscriptionListQueryBuilder
	if flt.NamespaceID != 0 {
		queryBuilder = queryBuilder.Where("namespace_id = ?", flt.NamespaceID)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.client.GetDB(ctx).QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptionsDomain []subscription.Subscription
	for rows.Next() {
		var subscriptionModel model.Subscription
		if err := rows.StructScan(&subscriptionModel); err != nil {
			return nil, err
		}

		subscriptionsDomain = append(subscriptionsDomain, *subscriptionModel.ToDomain())
	}

	return subscriptionsDomain, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {
	if sub == nil {
		return errors.New("subscription domain is nil")
	}

	subscriptionModel := new(model.Subscription)
	subscriptionModel.FromDomain(*sub)

	var newSubscriptionModel model.Subscription
	if err := r.client.GetDB(ctx).QueryRowxContext(ctx, subscriptionInsertQuery,
		subscriptionModel.NamespaceID,
		subscriptionModel.URN,
		subscriptionModel.Receiver,
		subscriptionModel.Match,
	).StructScan(&newSubscriptionModel); err != nil {
		err := checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return subscription.ErrRelation
		}
		return err
	}

	*sub = *newSubscriptionModel.ToDomain()

	return nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*subscription.Subscription, error) {
	query, args, err := subscriptionListQueryBuilder.Where("id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var subscriptionModel model.Subscription
	if err := r.client.GetDB(ctx).QueryRowxContext(ctx, query, args...).StructScan(&subscriptionModel); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, subscription.NotFoundError{ID: id}
		}
		return nil, err
	}

	return subscriptionModel.ToDomain(), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
	if sub == nil {
		return errors.New("subscription domain is nil")
	}

	subscriptionModel := new(model.Subscription)
	subscriptionModel.FromDomain(*sub)

	var newSubscriptionModel model.Subscription
	if err := r.client.GetDB(ctx).QueryRowxContext(ctx, subscriptionUpdateQuery,
		subscriptionModel.ID,
		subscriptionModel.NamespaceID,
		subscriptionModel.URN,
		subscriptionModel.Receiver,
		subscriptionModel.Match,
	).StructScan(&newSubscriptionModel); err != nil {
		err := checkPostgresError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return subscription.NotFoundError{ID: subscriptionModel.ID}
		}
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return subscription.ErrRelation
		}
		return err
	}

	*sub = *newSubscriptionModel.ToDomain()

	return nil
}

// TODO problem
func (r *SubscriptionRepository) Delete(ctx context.Context, id uint64) error {
	rows, err := r.client.GetDB(ctx).QueryxContext(ctx, subscriptionDeleteQuery, id)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (r *SubscriptionRepository) WithTransaction(ctx context.Context) context.Context {
	return r.client.WithTransaction(ctx, nil)
}

func (r *SubscriptionRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.client.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r *SubscriptionRepository) Commit(ctx context.Context) error {
	return r.client.Commit(ctx)
}
