package log

import (
	"context"
	"time"
)

type NotificationLogRepository interface {
	BulkCreate(context.Context, []Notification) error
	ListAlertIDsBySilenceID(context.Context, string) ([]int64, error)
	ListSubscriptionIDsBySilenceID(context.Context, string) ([]int64, error)
}

type Notification struct {
	ID             string
	NamespaceID    uint64
	NotificationID string
	SubscriptionID uint64
	ReceiverID     uint64
	AlertIDs       []int64
	SilenceIDs     []string
	CreatedAt      time.Time
}
