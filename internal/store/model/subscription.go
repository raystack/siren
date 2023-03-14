package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/pkg/pgc"
)

type SubscriptionReceiver struct {
	ID            uint64                 `json:"id"`
	Configuration map[string]interface{} `json:"configuration"`
}

type SubscriptionReceivers []SubscriptionReceiver

func (list *SubscriptionReceivers) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &list)
}

func (list SubscriptionReceivers) Value() (driver.Value, error) {
	val, err := json.Marshal(list)
	return string(val), err
}

type Subscription struct {
	ID          uint64                `db:"id"`
	NamespaceID uint64                `db:"namespace_id"`
	URN         string                `db:"urn"`
	Receiver    SubscriptionReceivers `db:"receiver"`
	Match       pgc.StringStringMap   `db:"match"`
	CreatedAt   time.Time             `db:"created_at"`
	UpdatedAt   time.Time             `db:"updated_at"`
}

func (s *Subscription) FromDomain(sub subscription.Subscription) {
	s.ID = sub.ID
	s.URN = sub.URN
	s.NamespaceID = sub.Namespace
	s.Match = sub.Match
	s.Receiver = make([]SubscriptionReceiver, 0)
	for _, item := range sub.Receivers {
		receiver := SubscriptionReceiver{
			ID:            item.ID,
			Configuration: item.Configuration,
		}
		s.Receiver = append(s.Receiver, receiver)
	}
	s.CreatedAt = sub.CreatedAt
	s.UpdatedAt = sub.UpdatedAt
}

func (s *Subscription) ToDomain() *subscription.Subscription {
	receivers := make([]subscription.Receiver, 0)
	for _, item := range s.Receiver {
		receiver := subscription.Receiver{
			ID:            item.ID,
			Configuration: item.Configuration,
		}
		receivers = append(receivers, receiver)
	}

	return &subscription.Subscription{
		ID:        s.ID,
		URN:       s.URN,
		Match:     s.Match,
		Namespace: s.NamespaceID,
		Receivers: receivers,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
