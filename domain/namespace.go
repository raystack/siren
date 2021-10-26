package domain

import "time"

type Namespace struct {
	Id          uint64                 `json:"id"`
	Urn         string                 `json:"urn"`
	Name        string                 `json:"name"`
	Provider    uint64                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type NamespaceService interface {
	ListNamespaces() ([]*Namespace, error)
	CreateNamespace(*Namespace) (*Namespace, error)
	GetNamespace(uint64) (*Namespace, error)
	UpdateNamespace(*Namespace) (*Namespace, error)
	DeleteNamespace(uint64) error
	Migrate() error
}
