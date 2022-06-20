package model

import (
	"time"

	"github.com/odpf/siren/core/namespace"
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

func (ns *Namespace) FromDomain(n *namespace.EncryptedNamespace) error {
	if n == nil {
		return nil
	}

	ns.Id = n.Id
	ns.Urn = n.Urn
	ns.Name = n.Name
	ns.ProviderId = n.Provider
	ns.Credentials = n.Credentials
	ns.Labels = StringStringMap(n.Labels)
	ns.CreatedAt = n.CreatedAt
	ns.UpdatedAt = n.UpdatedAt
	return nil
}

func (ns *Namespace) ToDomain() (*namespace.EncryptedNamespace, error) {
	if ns == nil {
		return nil, nil
	}

	return &namespace.EncryptedNamespace{
		Namespace: &namespace.Namespace{
			Id:        ns.Id,
			Urn:       ns.Urn,
			Name:      ns.Name,
			Provider:  ns.ProviderId,
			Labels:    ns.Labels,
			CreatedAt: ns.CreatedAt,
			UpdatedAt: ns.UpdatedAt,
		},
		Credentials: ns.Credentials,
	}, nil
}
