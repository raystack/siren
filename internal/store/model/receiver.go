package model

import (
	"time"

	"github.com/odpf/siren/core/receiver"
)

type Receiver struct {
	ID             uint64 `gorm:"primarykey"`
	Name           string
	Type           string
	Labels         StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	Configurations StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	Data           StringInterfaceMap `gorm:"-"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (rcv *Receiver) FromDomain(t *receiver.Receiver) {
	if t == nil {
		return
	}
	rcv.ID = t.ID
	rcv.Name = t.Name
	rcv.Type = t.Type
	rcv.Labels = t.Labels
	rcv.Configurations = t.Configurations
	rcv.Data = t.Data
	rcv.CreatedAt = t.CreatedAt
	rcv.UpdatedAt = t.UpdatedAt
}

func (rcv *Receiver) ToDomain() *receiver.Receiver {
	if rcv == nil {
		return nil
	}
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
