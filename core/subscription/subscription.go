package subscription

import (
	context "context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname SubscriptionRepository --filename subscription_repository.go --output=./mocks
type Repository interface {
	List(context.Context, Filter) ([]Subscription, error)
	Create(context.Context, *Subscription) error
	Get(context.Context, uint64) (*Subscription, error)
	Update(context.Context, *Subscription) error
	Delete(context.Context, uint64) error
}

type Receiver struct {
	ID            uint64                 `json:"id"`
	Configuration map[string]interface{} `json:"configuration"`

	// Type won't be exposed to the end-user, this is used to add more details for notification purposes
	Type string
}

type Subscription struct {
	ID        uint64            `json:"id"`
	URN       string            `json:"urn"`
	Namespace uint64            `json:"namespace"`
	Receivers []Receiver        `json:"receivers"`
	Match     map[string]string `json:"match"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
