package model

import (
	"github.com/odpf/siren/domain"
	"time"
)

type Receiver struct {
	Id             uint64 `gorm:"primarykey"`
	Name           string
	Type           string
	Labels         StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	Configurations StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	Data           StringInterfaceMap `gorm:"-"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (receiver *Receiver) FromDomain(t *domain.Receiver) *Receiver {
	if t == nil {
		return nil
	}
	receiver.Id = t.Id
	receiver.Name = t.Name
	receiver.Type = t.Type
	receiver.Labels = t.Labels
	receiver.Configurations = t.Configurations
	receiver.Data = t.Data
	receiver.CreatedAt = t.CreatedAt
	receiver.UpdatedAt = t.UpdatedAt
	return receiver
}

func (receiver *Receiver) ToDomain() *domain.Receiver {
	if receiver == nil {
		return nil
	}
	return &domain.Receiver{
		Id:             receiver.Id,
		Name:           receiver.Name,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: receiver.Configurations,
		Data:           receiver.Data,
		CreatedAt:      receiver.CreatedAt,
		UpdatedAt:      receiver.UpdatedAt,
	}
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackRepository interface {
	GetWorkspaceChannels(string) ([]Channel, error)
}

