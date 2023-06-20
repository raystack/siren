package model

import (
	"time"

	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/pkg/pgc"
)

type Receiver struct {
	ID             uint64              `db:"id"`
	Name           string              `db:"name"`
	Type           string              `db:"type"`
	Labels         pgc.StringStringMap `db:"labels"`
	Configurations pgc.StringAnyMap    `db:"configurations"`
	Data           pgc.StringAnyMap    `db:"-"` //TODO do we need this?
	CreatedAt      time.Time           `db:"created_at"`
	UpdatedAt      time.Time           `db:"updated_at"`
}

func (rcv *Receiver) FromDomain(t receiver.Receiver) {
	rcv.ID = t.ID
	rcv.Name = t.Name
	rcv.Type = t.Type
	rcv.Labels = t.Labels
	rcv.Configurations = pgc.StringAnyMap(t.Configurations)
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
