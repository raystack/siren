package namespace

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/provider"
	"github.com/pkg/errors"
	"time"
)

type StringStringMap map[string]string

func (m *StringStringMap) Scan(value interface{}) error {
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

type Namespace struct {
	Id          uint64 `gorm:"primarykey"`
	Provider    *provider.Provider
	ProviderId  uint64 `gorm:"uniqueIndex:name_provider_id_unique"`
	Name        string `gorm:"uniqueIndex:name_provider_id_unique"`
	Credentials string
	Labels      StringStringMap `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (namespace *Namespace) fromDomain(n *domain.Namespace) (*Namespace, error) {
	if n == nil {
		return nil, nil
	}
	namespace.Id = n.Id
	namespace.Name = n.Name
	namespace.ProviderId = n.Provider
	credentialsBytes, err := json.Marshal(n.Credentials)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	namespace.Credentials = string(credentialsBytes)
	namespace.Labels = n.Labels
	namespace.CreatedAt = n.CreatedAt
	namespace.UpdatedAt = n.UpdatedAt
	return namespace, nil
}

func (namespace *Namespace) toDomain() (*domain.Namespace, error) {
	if namespace == nil {
		return nil, nil
	}
	decryptedCredentials := make(map[string]interface{})
	err := json.Unmarshal([]byte(namespace.Credentials), &decryptedCredentials)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")

	}
	return &domain.Namespace{
		Id:          namespace.Id,
		Name:        namespace.Name,
		Provider:    namespace.ProviderId,
		Credentials: decryptedCredentials,
		Labels:      namespace.Labels,
		CreatedAt:   namespace.CreatedAt,
		UpdatedAt:   namespace.UpdatedAt,
	}, nil
}

type NamespaceRepository interface {
	Migrate() error
	List() ([]*Namespace, error)
	Create(*Namespace) (*Namespace, error)
	Get(uint64) (*Namespace, error)
	Update(*Namespace) (*Namespace, error)
	Delete(uint64) error
}
