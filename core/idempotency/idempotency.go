package idempotency

import (
	"time"
)

type Filter struct {
	TTL time.Duration
}

type Idempotency struct {
	ID        uint64
	Scope     string
	Key       string
	Success   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
