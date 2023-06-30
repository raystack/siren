package model

import (
	"time"

	"github.com/raystack/siren/core/provider"
	"github.com/raystack/siren/pkg/pgc"
)

type Provider struct {
	ID          uint64                 `db:"id"`
	Host        string                 `db:"host"`
	URN         string                 `db:"urn"`
	Name        string                 `db:"name"`
	Type        string                 `db:"type"`
	Credentials pgc.StringInterfaceMap `db:"credentials"`
	Labels      pgc.StringStringMap    `db:"labels"`
	CreatedAt   time.Time              `db:"created_at"`
	UpdatedAt   time.Time              `db:"updated_at"`
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
