package provider

import (
	"context"
	"time"
)

type Repository interface {
	List(context.Context, Filter) ([]Provider, error)
	Create(context.Context, *Provider) error
	Get(context.Context, uint64) (*Provider, error)
	Update(context.Context, *Provider) error
	Delete(context.Context, uint64) error
}

type Provider struct {
	ID          uint64            `json:"id"`
	URN         string            `json:"urn"`
	Host        string            `json:"host"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Credentials map[string]any    `json:"credentials"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
