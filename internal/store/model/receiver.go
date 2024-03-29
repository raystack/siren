package model

import (
	"time"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/pgc"
)

type Receiver struct {
	ID             uint64                 `db:"id"`
	Name           string                 `db:"name"`
	Type           string                 `db:"type"`
	Labels         pgc.StringStringMap    `db:"labels"`
	Configurations pgc.StringInterfaceMap `db:"configurations"`
	Data           pgc.StringInterfaceMap `db:"-"` //TODO do we need this?
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at"`
}

func (rcv *Receiver) FromDomain(t receiver.Receiver) {
	rcv.ID = t.ID
	rcv.Name = t.Name
	rcv.Type = t.Type
	rcv.Labels = t.Labels
	rcv.Configurations = pgc.StringInterfaceMap(t.Configurations)
	rcv.Data = t.Data
	rcv.CreatedAt = t.CreatedAt
	rcv.UpdatedAt = t.UpdatedAt
}

func (rcv *Receiver) ToDomain() *receiver.Receiver {
	return &receiver.Receiver{
		ID:             rcv.ID,
		Name:           rcv.Name,
		Type:           rcv.Type,
		Labels:         rcv.Labels,
		Configurations: rcv.Configurations,
		Data:           rcv.Data,
		CreatedAt:      rcv.CreatedAt,
		UpdatedAt:      rcv.UpdatedAt,
	}
}
