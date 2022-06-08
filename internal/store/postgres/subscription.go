package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store/model"
	"github.com/pkg/errors"
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

func (r *SubscriptionRepository) List(ctx context.Context) ([]*domain.Subscription, error) {
	var subscriptionModels []*model.Subscription
	selectQuery := "select * from subscriptions"
	result := r.getDb(ctx).Raw(selectQuery).Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []*domain.Subscription
	for _, s := range subscriptionModels {
		subscriptions = append(subscriptions, s.ToDomain())
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	m := new(model.Subscription)
	m.FromDomain(sub)
	if err := r.getDb(ctx).Create(m).Error; err != nil {
		return errors.Wrap(err, "failed to insert subscription")
	}

	newSubcription := m.ToDomain()
	*sub = *newSubcription
	return nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*domain.Subscription, error) {
	m := new(model.Subscription)
	result := r.getDb(ctx).Where(fmt.Sprintf("id = %d", id)).Find(m)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return m.ToDomain(), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	m := new(model.Subscription)
	m.FromDomain(sub)
	result := r.getDb(ctx).Where("id = ?", m.Id).Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("subscription doesn't exist")
	}

	if err := r.getDb(ctx).Where(fmt.Sprintf("id = %d", m.Id)).Find(m).Error; err != nil {
		return err
	}

	newSubcription := m.ToDomain()
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

func (r *SubscriptionRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Subscription{})
	if err != nil {
		return err
	}
	return nil
}
