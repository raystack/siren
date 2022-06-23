package provider

import (
	"context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname ProviderRepository --filename provider_repository.go --output=./mocks
type Repository interface {
	List(context.Context, Filter) ([]*Provider, error)
	Create(context.Context, *Provider) (uint64, error)
	Get(context.Context, uint64) (*Provider, error)
	Update(context.Context, *Provider) (uint64, error)
	Delete(context.Context, uint64) error
}

type Provider struct {
	ID          uint64                 `json:"id"`
	URN         string                 `json:"urn"`
	Host        string                 `json:"host"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Credentials map[string]interface{} `json:"credentials"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
