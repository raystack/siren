package model

import (
	"time"

	"github.com/goto/siren/core/notification"
)

type Idempotency struct {
	ID             uint64    `db:"id"`
	Scope          string    `db:"scope"`
	Key            string    `db:"key"`
	NotificationID string    `db:"notification_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (i *Idempotency) ToDomain() *notification.Idempotency {
	return &notification.Idempotency{
		ID:             i.ID,
		Scope:          i.Scope,
		Key:            i.Key,
		NotificationID: i.NotificationID,
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}
}
