package provider

import (
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname ProviderRepository --filename provider_repository.go --output=./mocks
type Repository interface {
	List(map[string]interface{}) ([]*Provider, error)
	Create(*Provider) (*Provider, error)
	Get(uint64) (*Provider, error)
	Update(*Provider) (*Provider, error)
	Delete(uint64) error
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
