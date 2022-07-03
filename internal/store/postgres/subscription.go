package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

// SubscriptionRepository talks to the store to read or insert data
type SubscriptionRepository struct {
	client *Client
}

// NewSubscriptionRepository returns SubscriptionRepository struct
func NewSubscriptionRepository(client *Client) *SubscriptionRepository {
	return &SubscriptionRepository{client}
}

func (r *SubscriptionRepository) List(ctx context.Context, flt subscription.Filter) ([]subscription.Subscription, error) {
	return r.list(ctx, r.client.db, flt)
}

func (r *SubscriptionRepository) list(ctx context.Context, tx *gorm.DB, flt subscription.Filter) ([]subscription.Subscription, error) {
	var subscriptionModels []*model.Subscription

	result := tx.WithContext(ctx)
	if flt.NamespaceID != 0 {
		result = result.Where("namespace_id = ?", flt.NamespaceID)
	}

	result = result.Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []subscription.Subscription
	for _, s := range subscriptionModels {
		subsDomain, err := s.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		subscriptions = append(subscriptions, *subsDomain)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) CreateWithTx(ctx context.Context, sub *subscription.Subscription, postProcessFn func(subs []subscription.Subscription) error) error {
	m := new(model.Subscription)
	if err := m.FromDomain(sub); err != nil {
		return err
	}

	if txErr := r.client.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&m).Error; err != nil {
			err = checkPostgresError(err)
			if errors.Is(err, errDuplicateKey) {
				return subscription.ErrDuplicate
			}
			if errors.Is(err, errForeignKeyViolation) {
				return subscription.ErrRelation
			}
			return err
		}

		// fetch all subscriptions in this namespace.
		subscriptionsInNamespace, err := r.list(ctx, tx, subscription.Filter{
			NamespaceID: sub.Namespace,
		})
		if err != nil {
			return err
		}

		return postProcessFn(subscriptionsInNamespace)
	}); txErr != nil {
		return txErr
	}

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

	subs, err := m.ToDomain()
	if err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *SubscriptionRepository) UpdateWithTx(ctx context.Context, sub *subscription.Subscription, postProcessFn func([]subscription.Subscription) error) error {
	m := new(model.Subscription)
	if err := m.FromDomain(sub); err != nil {
		return err
	}

	if txErr := r.client.db.Transaction(func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).Where("id = ?", m.ID).Updates(&m)
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

		// fetch all subscriptions in this namespace.
		subscriptionsInNamespace, err := r.list(ctx, tx, subscription.Filter{
			NamespaceID: sub.Namespace,
		})
		if err != nil {
			return err
		}

		return postProcessFn(subscriptionsInNamespace)
	}); txErr != nil {
		return txErr
	}

	return nil
}

func (r *SubscriptionRepository) DeleteWithTx(ctx context.Context, id uint64, namespaceID uint64, postProcessFn func([]subscription.Subscription) error) error {
	if txErr := r.client.db.Transaction(func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).Delete(model.Subscription{}, id)
		if result.Error != nil {
			return result.Error
		}
		// fetch all subscriptions in this namespace.
		subscriptionsInNamespace, err := r.list(ctx, tx, subscription.Filter{
			NamespaceID: namespaceID,
		})
		if err != nil {
			return err
		}
		return postProcessFn(subscriptionsInNamespace)
	}); txErr != nil {
		return txErr
	}
	return nil
}
