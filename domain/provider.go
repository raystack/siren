package domain

import "time"

type Provider struct {
	Id          uint64                 `json:"id"`
	Host        string                 `json:"host"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ProviderService interface {
	ListProviders() ([]*Provider, error)
	CreateProvider(*Provider) (*Provider, error)
	GetProvider(uint64) (*Provider, error)
	UpdateProvider(*Provider) (*Provider, error)
	DeleteProvider(uint64) error
	Migrate() error
}
