package domain

import "time"

type ReceiverMetadata struct {
	Id            uint64            `json:"id"`
	Configuration map[string]string `json:"configuration"`
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

type SubscriptionService interface {
	ListSubscriptions() ([]*Subscription, error)
	CreateSubscription(*Subscription) (*Subscription, error)
	GetSubscription(uint64) (*Subscription, error)
	UpdateSubscription(*Subscription) (*Subscription, error)
	DeleteSubscription(uint64) error
	Migrate() error
}
