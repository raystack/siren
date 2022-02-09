package model

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"time"
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

func (namespace *Namespace) FromDomain(n *domain.Namespace) (*Namespace, error) {
	if n == nil {
		return nil, nil
	}
	namespace.Id = n.Id
	namespace.Urn = n.Urn
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

func (namespace *Namespace) ToDomain() (*domain.Namespace, error) {
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
		Urn:         namespace.Urn,
		Name:        namespace.Name,
		Provider:    namespace.ProviderId,
		Credentials: decryptedCredentials,
		Labels:      namespace.Labels,
		CreatedAt:   namespace.CreatedAt,
		UpdatedAt:   namespace.UpdatedAt,
	}, nil
}

