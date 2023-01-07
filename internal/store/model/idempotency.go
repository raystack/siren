package model

import (
	"time"

	"github.com/odpf/siren/core/notification"
)

type Idempotency struct {
	ID        uint64    `db:"id"`
	Scope     string    `db:"scope"`
	Key       string    `db:"key"`
	Success   bool      `db:"success"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (i *Idempotency) ToDomain() *notification.Idempotency {
	return &notification.Idempotency{
		ID:        i.ID,
		Scope:     i.Scope,
		Key:       i.Key,
		Success:   i.Success,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}
