package notification

import (
	"context"
	"time"

	"github.com/goto/siren/pkg/errors"
)

const (
	ValidDurationRequestKey string = "valid_duration"

	FlowReceiver   string = "receiver"
	FlowSubscriber string = "subscriber"

	TypeAlert string = "alert"
	TypeEvent string = "event"
)

type Repository interface {
	Transactor
	Create(context.Context, Notification) (Notification, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

// Notification is a model of notification
type Notification struct {
	ID                string              `json:"id"`
	NamespaceID       uint64              `json:"namespace_id"`
	Type              string              `json:"type"`
	Data              map[string]any      `json:"data"`
	Labels            map[string]string   `json:"labels"`
	ValidDuration     time.Duration       `json:"valid_duration"`
	Template          string              `json:"template"`
	UniqueKey         string              `json:"unique_key"`
	ReceiverSelectors []map[string]string `json:"receiver_selectors"`
	CreatedAt         time.Time           `json:"created_at"`

	// won't be stored in notification table, only to propagate this to notification_subscriber
	AlertIDs []int64
}

func (n *Notification) EnrichID(id string) {
	if n == nil {
		return
	}
	n.ID = id

	if len(n.Data) == 0 {
		n.Data = map[string]any{}
	}

	n.Data["id"] = id
}

func (n Notification) Validate(flow string) error {
	if flow == FlowReceiver {
		if len(n.ReceiverSelectors) != 0 {
			return nil
		}
		return errors.ErrInvalid.WithCausef("notification type receiver should have receiver_selectors: %v", n)
	} else if flow == FlowSubscriber {
		if len(n.Labels) != 0 {
			return nil
		}
		return errors.ErrInvalid.WithCausef("notification type subscriber should have labels: %v", n)
	}

	return errors.ErrInvalid.WithCausef("invalid notification type: %v", n)
}
