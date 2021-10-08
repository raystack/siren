package domain

import (
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Receiver struct {
	Id             uint64                 `json:"id"`
	Urn            string                 `json:"urn"`
	Type           string                 `json:"type"`
	Labels         map[string]string      `json:"labels"`
	Configurations map[string]interface{} `json:"configurations"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ReceiverConfiguration interface {
	Get() (map[string]interface{}, error)
	Validate() error
}

type SlackConfigurations struct {
	Workspace string `json:"workspace" validate:"required"`
	Token     string `json:"token" validate:"required"`
}

func (sc *SlackConfigurations) Validate() error {
	v := validator.New()
	return v.Struct(sc)
}

type ReceiverService interface {
	ListReceivers() ([]*Receiver, error)
	CreateReceiver(*Receiver) (*Receiver, error)
	GetReceiver(uint64) (*Receiver, error)
	UpdateReceiver(*Receiver) (*Receiver, error)
	DeleteReceiver(uint64) error
	Migrate() error
}
