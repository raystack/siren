package receiver

import (
	"context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname ReceiverRepository --filename receiver_repository.go --output=./mocks
type Repository interface {
	List(context.Context, Filter) ([]Receiver, error)
	Create(context.Context, *Receiver) error
	Get(context.Context, uint64) (*Receiver, error)
	Update(context.Context, *Receiver) error
	Delete(context.Context, uint64) error
}

type Receiver struct {
	ID             uint64                 `json:"id"`
	Name           string                 `json:"name"`
	Labels         map[string]string      `json:"labels"`
	Configurations map[string]interface{} `json:"configurations"`
	Data           map[string]interface{} `json:"data"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`

	// Type should be immutable
	Type string `json:"type"`
}
