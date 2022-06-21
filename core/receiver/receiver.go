package receiver

import (
	"fmt"
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

type Configurations map[string]interface{}

func (c Configurations) GetString(key string) (string, error) {
	val, ok := c[key]
	if !ok {
		return "", fmt.Errorf("no value supplied for required configurations map key %q", key)
	}
	typedVal, ok := val.(string)
	if !ok {
		return "",
			fmt.Errorf(
				"wrong type for configurations map key %q: expected type %v, got value %q of type %t",
				key, "string", val, val)
	}
	return typedVal, nil
}

type Receiver struct {
	ID             uint64                 `json:"id"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Labels         map[string]string      `json:"labels"`
	Configurations Configurations         `json:"configurations"`
	Data           map[string]interface{} `json:"data"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}
