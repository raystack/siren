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
	*transaction
}

// NewSubscriptionRepository returns SubscriptionRepository struct
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{&transaction{db}}
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]*subscription.Subscription, error) {
	var subscriptionModels []*model.Subscription
	selectQuery := "select * from subscriptions"
	result := r.getDb(ctx).Raw(selectQuery).Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []*subscription.Subscription
	for _, s := range subscriptionModels {
		subsDomain, err := s.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		subscriptions = append(subscriptions, subsDomain)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *subscription.Subscription) error {
	m := new(model.Subscription)
	if err := m.FromDomain(sub); err != nil {
		return err
	}
	if err := r.getDb(ctx).Create(m).Error; err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		return err
	}

	newSubcription, err := m.ToDomain()
	if err != nil {
		return err
	}
	*sub = *newSubcription
	return nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*subscription.Subscription, error) {
	m := new(model.Subscription)
	result := r.getDb(ctx).Where(fmt.Sprintf("id = %d", id)).Find(m)
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

func (r *SubscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
	m := new(model.Subscription)
	if err := m.FromDomain(sub); err != nil {
		return err
	}
	result := r.getDb(ctx).Where("id = ?", m.ID).Updates(m)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return subscription.ErrDuplicate
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return subscription.NotFoundError{ID: sub.ID}
	}

	if err := r.getDb(ctx).Where(fmt.Sprintf("id = %d", m.ID)).Find(m).Error; err != nil {
		return err
	}

	newSubcription, err := m.ToDomain()
	if err != nil {
		return errors.New("failed to convert subscription from model to domain")
	}
	*sub = *newSubcription
	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uint64) error {
	result := r.getDb(ctx).Delete(model.Subscription{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
