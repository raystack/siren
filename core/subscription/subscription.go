package subscription

import (
	context "context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname SubscriptionRepository --filename subscription_repository.go --output=./mocks
type Repository interface {
	Transactor
	Migrate() error
	List(context.Context) ([]*Subscription, error)
	Create(context.Context, *Subscription) error
	Get(context.Context, uint64) (*Subscription, error)
	Update(context.Context, *Subscription) error
	Delete(context.Context, uint64) error
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type Subscription struct {
	Id        uint64             `json:"id"`
	Urn       string             `json:"urn"`
	Namespace uint64             `json:"namespace"`
	Receivers []ReceiverMetadata `json:"receivers"`
	Match     map[string]string  `json:"match"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
