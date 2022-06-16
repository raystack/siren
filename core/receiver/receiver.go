package receiver

import (
	"time"
)

const (
	TypeSlack     string = "slack"
	TypeHTTP      string = "http"
	TypePagerDuty string = "pagerduty"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname ReceiverRepository --filename receiver_repository.go --output=./mocks
type Repository interface {
	Migrate() error
	List() ([]*Receiver, error)
	Create(*Receiver) error
	Get(uint64) (*Receiver, error)
	Update(*Receiver) error
	Delete(uint64) error
}

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
