package domain

import (
	"time"
)

type Receiver struct {
	Id             uint64                 `json:"id"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Labels         map[string]string      `json:"labels"`
	Configurations map[string]interface{} `json:"configurations"`
	Data           map[string]interface{} `json:"data"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}


type ReceiverService interface {
	ListReceivers() ([]*Receiver, error)
	CreateReceiver(*Receiver) (*Receiver, error)
	GetReceiver(uint64) (*Receiver, error)
	UpdateReceiver(*Receiver) (*Receiver, error)
	DeleteReceiver(uint64) error
	Migrate() error
}
