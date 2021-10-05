package provider

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

type Provider struct {
	Id          uint64 `gorm:"primarykey"`
	Host        string
	Name        string
	Type        string
	Credentials StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	Labels      StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Provider *Provider) fromDomain(t *domain.Provider) *Provider {
	if t == nil {
		return nil
	}
	Provider.Id = t.Id
	Provider.Host = t.Host
	Provider.Name = t.Name
	Provider.Type = t.Type
	Provider.Credentials = t.Credentials
	Provider.Labels = t.Labels
	Provider.CreatedAt = t.CreatedAt
	Provider.UpdatedAt = t.UpdatedAt
	return Provider
}

func (Provider *Provider) toDomain() *domain.Provider {
	if Provider == nil {
		return nil
	}
	return &domain.Provider{
		Id:          Provider.Id,
		Host:        Provider.Host,
		Name:        Provider.Name,
		Type:        Provider.Type,
		Credentials: Provider.Credentials,
		Labels:      Provider.Labels,
		CreatedAt:   Provider.CreatedAt,
		UpdatedAt:   Provider.UpdatedAt,
	}
}

type ProviderRepository interface {
	Migrate() error
	List() ([]*Provider, error)
	Create(*Provider) (*Provider, error)
	Get(uint64) (*Provider, error)
	Update(*Provider) (*Provider, error)
	Delete(uint64) error
}
