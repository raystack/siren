package model

import (
	"time"

	"github.com/odpf/siren/domain"
)

type Namespace struct {
	Id          uint64 `gorm:"primarykey"`
	Provider    *Provider
	ProviderId  uint64 `gorm:"uniqueIndex:urn_provider_id_unique"`
	Urn         string `gorm:"uniqueIndex:urn_provider_id_unique"`
	Name        string
	Credentials string
	Labels      StringStringMap `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (namespace *Namespace) FromDomain(n *domain.EncryptedNamespace) error {
	if n == nil {
		return nil
	}

	namespace.Id = n.Id
	namespace.Urn = n.Urn
	namespace.Name = n.Name
	namespace.ProviderId = n.Provider
	namespace.Credentials = n.Credentials
	namespace.Labels = StringStringMap(n.Labels)
	namespace.CreatedAt = n.CreatedAt
	namespace.UpdatedAt = n.UpdatedAt
	return nil
}

func (namespace *Namespace) ToDomain() (*domain.EncryptedNamespace, error) {
	if namespace == nil {
		return nil, nil
	}

	return &domain.EncryptedNamespace{
		Namespace: &domain.Namespace{
			Id:        namespace.Id,
			Urn:       namespace.Urn,
			Name:      namespace.Name,
			Provider:  namespace.ProviderId,
			Labels:    namespace.Labels,
			CreatedAt: namespace.CreatedAt,
			UpdatedAt: namespace.UpdatedAt,
		},
		Credentials: namespace.Credentials,
	}, nil
}
