package notification

import (
	"context"
	"strconv"
	"time"

	"github.com/goto/siren/pkg/errors"
)

const (
	ReceiverIDLabelKey      string = "receiver_id"
	ValidDurationRequestKey string = "valid_duration"

	TypeReceiver   string = "receiver"
	TypeSubscriber string = "subscriber"
)

type Repository interface {
	Create(context.Context, Notification) (Notification, error)
}

// Notification is a model of notification
// if type is `receiver`, it is expected for the labels to have
// receiver_id = int
type Notification struct {
	ID            string            `json:"id"`
	NamespaceID   uint64            `json:"namespace_id"`
	Type          string            `json:"type"`
	Data          map[string]any    `json:"data"`
	Labels        map[string]string `json:"labels"`
	ValidDuration time.Duration     `json:"valid_duration"`
	Template      string            `json:"template"`
	UniqueKey     string            `json:"unique_key"`
	CreatedAt     time.Time         `json:"created_at"`

	// won't be stored in notification table, only to propaget this to notification_subscriber
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

func (n Notification) Validate() error {
	if n.Type == TypeReceiver {
		if v, ok := n.Labels[ReceiverIDLabelKey]; ok {
			intVar, err := strconv.ParseInt(v, 0, 64)
			if err == nil && intVar != 0 {
				return nil
			}
		}
		return errors.ErrInvalid.WithCausef("notification type receiver should have valid receiver_id: %v", n)
	} else if n.Type == TypeSubscriber {
		if len(n.Labels) != 0 {
			return nil
		}
		return errors.ErrInvalid.WithCausef("notification type subscriber should have labels: %v", n)
	}

	return errors.ErrInvalid.WithCausef("invalid notification type: %v", n)
}
