package notification

import (
	"context"
	"time"
)

type IdempotencyFilter struct {
	TTL time.Duration
}

//go:generate mockery --name=IdempotencyRepository -r --case underscore --with-expecter --structname IdempotencyRepository --filename idempotency_repository.go --output=./mocks
type IdempotencyRepository interface {
	InsertOnConflictReturning(context.Context, string, string) (*Idempotency, error)
	UpdateSuccess(context.Context, uint64, bool) error
	Delete(context.Context, IdempotencyFilter) error
}

type Idempotency struct {
	ID        uint64
	Scope     string
	Key       string
	Success   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
