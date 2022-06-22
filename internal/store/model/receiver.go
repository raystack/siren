package model

import (
	"time"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/pkg/errors"
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

func (rcv *Receiver) FromDomain(t *receiver.Receiver) error {
	if t == nil {
		return errors.New("receiver domain is nil")
	}
	rcv.ID = t.ID
	rcv.Name = t.Name
	rcv.Type = t.Type
	rcv.Labels = t.Labels
	rcv.Configurations = StringInterfaceMap(t.Configurations)
	rcv.Data = t.Data
	rcv.CreatedAt = t.CreatedAt
	rcv.UpdatedAt = t.UpdatedAt
	return nil
}

func (rcv *Receiver) ToDomain() (*receiver.Receiver, error) {
	if rcv == nil {
		return nil, errors.New("receiver model is nil")
	}
	return &receiver.Receiver{
		ID:             rcv.ID,
		Name:           rcv.Name,
		Type:           rcv.Type,
		Labels:         rcv.Labels,
		Configurations: receiver.Configurations(rcv.Configurations),
		Data:           rcv.Data,
		CreatedAt:      rcv.CreatedAt,
		UpdatedAt:      rcv.UpdatedAt,
	}, nil
}
