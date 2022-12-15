package jobs

import (
	"context"
	"time"

	"github.com/odpf/salt/log"
)

//go:generate mockery --name=NotificationService -r --case underscore --with-expecter --structname NotificationService --filename notification_service.go --output=./mocks
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
