package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/pkg/errors"
)

type ReceiverMetadata struct {
	ID            uint64            `json:"id"`
	Configuration map[string]string `json:"configuration"`
}

type ReceiverMetadataList []ReceiverMetadata

func (list *ReceiverMetadataList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &list)
}

func (list ReceiverMetadataList) Value() (driver.Value, error) {
	val, err := json.Marshal(list)
	return string(val), err
}

type Subscription struct {
	ID          uint64 `gorm:"primarykey"`
	Namespace   *Namespace
	NamespaceId uint64
	URN         string               `gorm:"unique"`
	Receiver    ReceiverMetadataList `gorm:"type:jsonb" sql:"type:jsonb" `
	Match       StringStringMap      `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Subscription) FromDomain(sub *subscription.Subscription) error {
	if s == nil {
		return errors.New("subscription domain is nil")
	}
	s.ID = sub.ID
	s.URN = sub.URN
	s.NamespaceId = sub.Namespace
	s.Match = sub.Match
	s.Receiver = make([]ReceiverMetadata, 0)
	for _, item := range sub.Receivers {
		receiver := ReceiverMetadata{
			ID:            item.ID,
			Configuration: item.Configuration,
		}
		s.Receiver = append(s.Receiver, receiver)
	}
	s.CreatedAt = sub.CreatedAt
	s.UpdatedAt = sub.UpdatedAt
	return nil
}

func (s *Subscription) ToDomain() (*subscription.Subscription, error) {
	if s == nil {
		return nil, errors.New("subscription model is nil")
	}
	receivers := make([]subscription.ReceiverMetadata, 0)
	for _, item := range s.Receiver {
		receiver := subscription.ReceiverMetadata{
			ID:            item.ID,
			Configuration: item.Configuration,
		}
		receivers = append(receivers, receiver)
	}

	return &subscription.Subscription{
		ID:        s.ID,
		URN:       s.URN,
		Match:     s.Match,
		Namespace: s.NamespaceId,
		Receivers: receivers,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}, nil
}
