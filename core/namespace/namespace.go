package namespace

import (
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname NamespaceRepository --filename namespace_repository.go --output=./mocks
type Repository interface {
	List() ([]*EncryptedNamespace, error)
	Create(*EncryptedNamespace) error
	Get(uint64) (*EncryptedNamespace, error)
	Update(*EncryptedNamespace) error
	Delete(uint64) error
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
