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
	Urn         string `gorm:"unique:urn"`
	Name        string
	Type        string
	Credentials StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	Labels      StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (provider *Provider) fromDomain(t *domain.Provider) *Provider {
	if t == nil {
		return nil
	}
	provider.Id = t.Id
	provider.Host = t.Host
	provider.Urn = t.Urn
	provider.Name = t.Name
	provider.Type = t.Type
	provider.Credentials = t.Credentials
	provider.Labels = t.Labels
	provider.CreatedAt = t.CreatedAt
	provider.UpdatedAt = t.UpdatedAt
	return provider
}

func (provider *Provider) toDomain() *domain.Provider {
	if provider == nil {
		return nil
	}
	return &domain.Provider{
		Id:          provider.Id,
		Host:        provider.Host,
		Name:        provider.Name,
		Urn:         provider.Urn,
		Type:        provider.Type,
		Credentials: provider.Credentials,
		Labels:      provider.Labels,
		CreatedAt:   provider.CreatedAt,
		UpdatedAt:   provider.UpdatedAt,
	}
}

type ProviderRepository interface {
	Migrate() error
	List(map[string]interface{}) ([]*Provider, error)
	Create(*Provider) (*Provider, error)
	Get(uint64) (*Provider, error)
	Update(*Provider) (*Provider, error)
	Delete(uint64) error
}
