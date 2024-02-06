package jobs

import (
	"context"
	"time"

	"github.com/goto/salt/log"
)

type NotificationService interface {
	RemoveIdempotencies(ctx context.Context, TTL time.Duration) error
}

type handler struct {
	logger              log.Logger
	notificationService NotificationService
}

func NewHandler(
	logger log.Logger,
	notificationService NotificationService,
) *handler {
	return &handler{
		logger:              logger,
		notificationService: notificationService,
	}
}
