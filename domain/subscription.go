package domain

import (
	"context"
	"time"
)

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
	ListSubscriptions(context.Context) ([]*Subscription, error)
	CreateSubscription(context.Context, *Subscription) error
	GetSubscription(context.Context, uint64) (*Subscription, error)
	UpdateSubscription(context.Context, *Subscription) error
	DeleteSubscription(context.Context, uint64) error
	Migrate() error
}
