package namespace

import (
	"context"
	"time"

	"github.com/raystack/siren/core/provider"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname NamespaceRepository --filename namespace_repository.go --output=./mocks
type Repository interface {
	Transactor
	List(context.Context) ([]EncryptedNamespace, error)
	Create(context.Context, *EncryptedNamespace) error
	Get(context.Context, uint64) (*EncryptedNamespace, error)
	Update(context.Context, *EncryptedNamespace) error
	Delete(context.Context, uint64) error
}

//go:generate mockery --name=Transactor -r --case underscore --with-expecter --structname Transactor --filename transactor.go --output=./mocks
type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type EncryptedNamespace struct {
	*Namespace
	CredentialString string
}

type Namespace struct {
	ID          uint64                 `json:"id"`
	URN         string                 `json:"urn"`
	Name        string                 `json:"name"`
	Provider    provider.Provider      `json:"provider"`
	Credentials map[string]interface{} `json:"credentials"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}
