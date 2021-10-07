package domain

import "time"

type Receiver struct {
	Id            uint64            `json:"id"`
	Urn           string            `json:"urn"`
	Type          string            `json:"type"`
	Labels        map[string]string `json:"labels"`
	Configuration map[string]string `json:"configuration"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type ReceiverService interface {
	ListReceivers() ([]*Receiver, error)
	CreateReceiver(*Receiver) (*Receiver, error)
	GetReceiver(uint64) (*Receiver, error)
	UpdateReceiver(*Receiver) (*Receiver, error)
	DeleteReceiver(uint64) error
	Migrate() error
}
