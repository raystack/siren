package notification

import (
	"context"
	"time"
)

type IdempotencyFilter struct {
	TTL time.Duration
}

type IdempotencyRepository interface {
	Create(ctx context.Context, scope, key, notificationID string) (*Idempotency, error)
	Check(ctx context.Context, scope, key string) (*Idempotency, error)
	Delete(context.Context, IdempotencyFilter) error
}

type Idempotency struct {
	ID             uint64    `json:"id"`
	Scope          string    `json:"scope"`
	Key            string    `json:"key"`
	NotificationID string    `json:"notification_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
