package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/pkg/errors"
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
	ID          uint64 `gorm:"primarykey"`
	Host        string
	URN         string `gorm:"unique:urn"`
	Name        string
	Type        string
	Credentials StringInterfaceMap `gorm:"type:jsonb" sql:"type:jsonb" `
	Labels      StringStringMap    `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Provider) FromDomain(t *provider.Provider) error {
	if t == nil {
		return errors.New("provider domain is nil")
	}
	p.ID = t.ID
	p.Host = t.Host
	p.URN = t.URN
	p.Name = t.Name
	p.Type = t.Type
	p.Credentials = t.Credentials
	p.Labels = t.Labels
	p.CreatedAt = t.CreatedAt
	p.UpdatedAt = t.UpdatedAt
	return nil
}

func (p *Provider) ToDomain() (*provider.Provider, error) {
	if p == nil {
		return nil, errors.New("provider model is nil")
	}
	return &provider.Provider{
		ID:          p.ID,
		Host:        p.Host,
		Name:        p.Name,
		URN:         p.URN,
		Type:        p.Type,
		Credentials: p.Credentials,
		Labels:      p.Labels,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}
