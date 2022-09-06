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
	if value == nil {
		m = new(StringInterfaceMap)
		return nil
	}
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
	if value == nil {
		m = new(StringStringMap)
		return nil
	}
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
	ID          uint64             `db:"id"`
	Host        string             `db:"host"`
	URN         string             `db:"urn"`
	Name        string             `db:"name"`
	Type        string             `db:"type"`
	Credentials StringInterfaceMap `db:"credentials"`
	Labels      StringStringMap    `db:"labels"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
}

func (p *Provider) FromDomain(t provider.Provider) {
	p.ID = t.ID
	p.Host = t.Host
	p.URN = t.URN
	p.Name = t.Name
	p.Type = t.Type
	p.Credentials = t.Credentials
	p.Labels = t.Labels
	p.CreatedAt = t.CreatedAt
	p.UpdatedAt = t.UpdatedAt
}

func (p *Provider) ToDomain() *provider.Provider {
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
	}
}
