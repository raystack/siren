package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

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
	var subscriptionModels []*model.Subscription

	result := r.client.GetDB(ctx).WithContext(ctx)
	if flt.NamespaceID != 0 {
		result = result.Where("namespace_id = ?", flt.NamespaceID)
	}

	result = result.Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []subscription.Subscription
	for _, s := range subscriptionModels {
		subscriptions = append(subscriptions, *s.ToDomain())
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {
	if sub == nil {
		return errors.New("subscription domain is nil")
	}

	m := new(model.Subscription)
	m.FromDomain(sub)

	if err := r.client.GetDB(ctx).WithContext(ctx).Create(&m).Error; err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return subscription.ErrRelation
		}
		return err
	}

	*sub = *m.ToDomain()

	return nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*subscription.Subscription, error) {
	m := new(model.Subscription)

	result := r.client.db.WithContext(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&m)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, subscription.NotFoundError{ID: id}
	}

	return m.ToDomain(), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
	if sub == nil {
		return errors.New("subscription domain is nil")
	}

	m := new(model.Subscription)
	m.FromDomain(sub)

	result := r.client.GetDB(ctx).WithContext(ctx).Where("id = ?", m.ID).Updates(&m)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return subscription.ErrRelation
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return subscription.NotFoundError{ID: m.ID}
	}

	*sub = *m.ToDomain()

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uint64, namespaceID uint64) error {
	result := r.client.GetDB(ctx).WithContext(ctx).Delete(model.Subscription{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *SubscriptionRepository) WithTransaction(ctx context.Context) context.Context {
	return r.client.WithTransaction(ctx)
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
