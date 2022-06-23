package namespace

import (
	"context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname NamespaceRepository --filename namespace_repository.go --output=./mocks
type Repository interface {
	List(context.Context) ([]*EncryptedNamespace, error)
	Create(context.Context, *EncryptedNamespace) (uint64, error)
	Get(context.Context, uint64) (*EncryptedNamespace, error)
	Update(context.Context, *EncryptedNamespace) (uint64, error)
	Delete(context.Context, uint64) error
}

type Namespace struct {
	ID          uint64                 `json:"id"`
	URN         string                 `json:"urn"`
	Name        string                 `json:"name"`
	Provider    uint64                 `json:"provider"`
	Credentials map[string]interface{} `json:"credentials"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type EncryptedNamespace struct {
	*Namespace
	Credentials string
}
