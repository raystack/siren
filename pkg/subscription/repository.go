package subscription

import (
	"context"
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/subscription/alertmanager"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	transactionContextKey = struct{}{}
)

// Repository talks to the store to read or insert data
type Repository struct {
	db       *gorm.DB
	amClient alertmanager.Client
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db, nil}
}

func (r *Repository) WithTransaction(ctx context.Context) context.Context {
	tx := r.db.Begin()
	return context.WithValue(ctx, transactionContextKey, tx)
}

func (r *Repository) Rollback(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Rollback()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (r *Repository) Commit(ctx context.Context) error {
	if tx := extractTransaction(ctx); tx != nil {
		tx = tx.Commit()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (r *Repository) List(ctx context.Context) ([]*domain.Subscription, error) {
	var subscriptionModels []*Subscription
	selectQuery := "select * from subscriptions"
	result := r.getDb(ctx).Raw(selectQuery).Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []*domain.Subscription
	for _, s := range subscriptionModels {
		subscriptions = append(subscriptions, s.toDomain())
	}

	return subscriptions, nil
}

func (r *Repository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	var model Subscription
	model.fromDomain(sub)
	var newSubscription *Subscription
	var err error

	newSubscription, err = r.insertSubscriptionIntoDB(r.getDb(ctx), &model)
	if err != nil {
		return nil, errors.Wrap(err, "r.insertSubscriptionIntoDB")
	}

	return newSubscription.toDomain(), nil
}

func (r *Repository) Get(ctx context.Context, id uint64) (*domain.Subscription, error) {
	var subscription Subscription
	result := r.getDb(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return subscription.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	model := new(Subscription)
	model.fromDomain(sub)
	result := r.getDb(ctx).Where("id = ?", model.Id).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("subscription doesn't exist")
	}
	return model.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	result := r.getDb(ctx).Delete(Subscription{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) Migrate() error {
	err := r.db.AutoMigrate(&Subscription{})
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) insertSubscriptionIntoDB(tx *gorm.DB, sub *Subscription) (*Subscription, error) {
	var newSubscription Subscription
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

func (r *Repository) getDb(ctx context.Context) *gorm.DB {
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
