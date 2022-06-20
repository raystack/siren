package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/odpf/siren/core/provider"
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

func (p *Provider) FromDomain(t *provider.Provider) *Provider {
	if t == nil {
		return nil
	}
	p.Id = t.Id
	p.Host = t.Host
	p.Urn = t.Urn
	p.Name = t.Name
	p.Type = t.Type
	p.Credentials = t.Credentials
	p.Labels = t.Labels
	p.CreatedAt = t.CreatedAt
	p.UpdatedAt = t.UpdatedAt
	return p
}

func (p *Provider) ToDomain() *provider.Provider {
	if p == nil {
		return nil
	}
	return &provider.Provider{
		Id:          p.Id,
		Host:        p.Host,
		Name:        p.Name,
		Urn:         p.Urn,
		Type:        p.Type,
		Credentials: p.Credentials,
		Labels:      p.Labels,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
