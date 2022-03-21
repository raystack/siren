package subscription

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
	"github.com/pkg/errors"
)

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

func (list *ReceiverMetadataList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &list)
}

func (list ReceiverMetadataList) Value() (driver.Value, error) {
	val, err := json.Marshal(list)
	return string(val), err
}

type ReceiverMetadata struct {
	Id            uint64            `json:"id"`
	Configuration map[string]string `json:"configuration"`
}

type EnrichedReceiverMetadata struct {
	Id            uint64            `json:"id"`
	Type          string            `json:"type"`
	Configuration map[string]string `json:"configuration"`
}

type ReceiverMetadataList []ReceiverMetadata
type EnrichedReceiverMetadataList []EnrichedReceiverMetadata

type Subscription struct {
	Id          uint64 `gorm:"primarykey"`
	Namespace   *model.Namespace
	NamespaceId uint64
	Urn         string               `gorm:"unique"`
	Receiver    ReceiverMetadataList `gorm:"type:jsonb" sql:"type:jsonb" `
	Match       StringStringMap      `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SubscriptionEnrichedWithReceivers struct {
	Id          uint64
	Namespace   *model.Namespace
	NamespaceId uint64
	Urn         string
	Receiver    EnrichedReceiverMetadataList
	Match       StringStringMap
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Subscription) fromDomain(sub *domain.Subscription) *Subscription {
	if s == nil {
		return nil
	}
	s.Id = sub.Id
	s.Urn = sub.Urn
	s.NamespaceId = sub.Namespace
	s.Match = sub.Match
	s.Receiver = make([]ReceiverMetadata, 0)
	for _, item := range sub.Receivers {
		receiver := ReceiverMetadata{
			Id:            item.Id,
			Configuration: item.Configuration,
		}
		s.Receiver = append(s.Receiver, receiver)
	}
	s.CreatedAt = sub.CreatedAt
	s.UpdatedAt = sub.UpdatedAt
	return s
}

func (s *Subscription) toDomain() *domain.Subscription {
	if s == nil {
		return nil
	}
	receivers := make([]domain.ReceiverMetadata, 0)
	for _, item := range s.Receiver {
		receiver := domain.ReceiverMetadata{
			Id:            item.Id,
			Configuration: item.Configuration,
		}
		receivers = append(receivers, receiver)
	}
	subscription := &domain.Subscription{
		Id:        s.Id,
		Urn:       s.Urn,
		Match:     s.Match,
		Namespace: s.NamespaceId,
		Receivers: receivers,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
	return subscription
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type SubscriptionRepository interface {
	Transactor
	Migrate() error
	List(context.Context) ([]*domain.Subscription, error)
	Create(context.Context, *domain.Subscription) (*domain.Subscription, error)
	Get(context.Context, uint64) (*domain.Subscription, error)
	Update(context.Context, *domain.Subscription) (*domain.Subscription, error)
	Delete(context.Context, uint64, domain.NamespaceService, domain.ProviderService, domain.ReceiverService) error
}
