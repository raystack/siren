package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	transactionContextKey = struct{}{}
)

// SubscriptionRepository talks to the store to read or insert data
type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository returns SubscriptionRepository struct
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

func (r *SubscriptionRepository) WithTransaction(ctx context.Context) context.Context {
	tx := r.db.Begin()
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (r *SubscriptionRepository) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Rollback()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (r *SubscriptionRepository) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Commit()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
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

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	var m model.Subscription
	m.FromDomain(sub)
	var newSubscription *model.Subscription
	var err error

	newSubscription, err = r.insertSubscriptionIntoDB(r.getDb(ctx), &m)
	if err != nil {
		return nil, errors.Wrap(err, "r.insertSubscriptionIntoDB")
	}

	return newSubscription.ToDomain(), nil
}

func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*domain.Subscription, error) {
	var subscription model.Subscription
	result := r.getDb(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return subscription.ToDomain(), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	m := new(model.Subscription)
	m.FromDomain(sub)
	result := r.getDb(ctx).Where("id = ?", m.Id).Updates(m)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("subscription doesn't exist")
	}
	return m.ToDomain(), nil
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

func (r *SubscriptionRepository) insertSubscriptionIntoDB(tx *gorm.DB, sub *model.Subscription) (*model.Subscription, error) {
	var newSubscription model.Subscription
	result := tx.Create(sub)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to insert subscription")
	}

	result = tx.Where(fmt.Sprintf("id = %d", sub.Id)).Find(&newSubscription)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get newly inserted subscription")
	}
	return &newSubscription, nil
}

func (r *SubscriptionRepository) getDb(ctx context.Context) *gorm.DB {
	db := r.db
	if tx := extractTransaction(ctx); tx != nil {
		db = tx
	}
	return db
}

func extractTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); !ok {
		return nil
	} else {
		return tx
	}
}
