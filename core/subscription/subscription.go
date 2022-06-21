package subscription

import (
	context "context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/pkg/errors"
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
	ID        uint64             `json:"id"`
	URN       string             `json:"urn"`
	Namespace uint64             `json:"namespace"`
	Receivers []ReceiverMetadata `json:"receivers"`
	Match     map[string]string  `json:"match"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type StringStringMap map[string]string

func (m *StringStringMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringStringMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

type ReceiverMetadata struct {
	ID            uint64            `json:"id"`
	Configuration map[string]string `json:"configuration"`
}

type EnrichedReceiverMetadata struct {
	ID            uint64            `json:"id"`
	Type          string            `json:"type"`
	Configuration map[string]string `json:"configuration"`
}

type EnrichedReceiverMetadataList []EnrichedReceiverMetadata

type SubscriptionEnrichedWithReceivers struct {
	ID          uint64
	Namespace   *namespace.Namespace
	NamespaceId uint64
	URN         string
	Receiver    EnrichedReceiverMetadataList
	Match       StringStringMap
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
