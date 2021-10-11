package receiver

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/odpf/siren/domain"
	"time"
)

type StringInterfaceMap map[string]interface{}
type StringStringMap map[string]string

func (m *StringInterfaceMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringInterfaceMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

func (m *StringStringMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("failed type assertion to []byte")
	}
	return json.Unmarshal(b, &m)
}

func (a StringStringMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}

type Receiver struct {
	Id             uint64 `gorm:"primarykey"`
	Urn            string
	Type           string
	Labels         StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	Configurations StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (receiver *Receiver) fromDomain(t *domain.Receiver) *Receiver {
	if t == nil {
		return nil
	}
	receiver.Id = t.Id
	receiver.Urn = t.Urn
	receiver.Type = t.Type
	receiver.Labels = t.Labels
	receiver.Configurations = t.Configurations
	receiver.CreatedAt = t.CreatedAt
	receiver.UpdatedAt = t.UpdatedAt
	return receiver
}

func (receiver *Receiver) toDomain() *domain.Receiver {
	if receiver == nil {
		return nil
	}
	return &domain.Receiver{
		Id:             receiver.Id,
		Urn:            receiver.Urn,
		Type:           receiver.Type,
		Labels:         receiver.Labels,
		Configurations: receiver.Configurations,
		CreatedAt:      receiver.CreatedAt,
		UpdatedAt:      receiver.UpdatedAt,
	}
}

type ReceiverRepository interface {
	Migrate() error
	List() ([]*Receiver, error)
	Create(*Receiver) (*Receiver, error)
	Get(uint64) (*Receiver, error)
	Update(*Receiver) (*Receiver, error)
	Delete(uint64) error
}
