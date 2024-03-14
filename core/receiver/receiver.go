package receiver

import (
	"context"
	"fmt"
	"time"
)

const (
	LabelKeyID   = "id"
	LabelKeyType = "type"
)

type Repository interface {
	List(context.Context, Filter) ([]Receiver, error)
	Create(context.Context, *Receiver) error
	Get(context.Context, uint64, Filter) (*Receiver, error)
	Update(context.Context, *Receiver) error
	Delete(context.Context, uint64) error
	PatchLabels(context.Context, *Receiver) error
}

type Receiver struct {
	ID             uint64            `json:"id"`
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels"`
	Configurations map[string]any    `json:"configurations"`
	Data           map[string]any    `json:"data"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`

	// The pointer to receiver parent of a child receiver. This field is required if a receiver is a child receiver
	// If ParentID != 0, the receiver is a child receiver.
	ParentID uint64 `json:"parent_id"`

	// Type should be immutable
	Type string `json:"type"`
}

func (r *Receiver) Validate() error {
	if r.Type == TypeSlackChannel && r.ParentID == 0 {
		return fmt.Errorf("type slack_channel needs receiver parent ID")
	}

	return nil
}

func (r *Receiver) enrichPredefinedLabels() {
	// populate predefined labels: id, type
	if len(r.Labels) == 0 {
		r.Labels = make(map[string]string)
	}
	r.Labels[LabelKeyID] = fmt.Sprintf("%d", r.ID)
	r.Labels[LabelKeyType] = r.Type
}
