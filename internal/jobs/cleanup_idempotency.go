package jobs

import (
	"context"
	"time"
)

func (h *handler) CleanupIdempotencies(ctx context.Context, TTL time.Duration) error {
	return h.notificationService.RemoveIdempotencies(ctx, TTL)
}
