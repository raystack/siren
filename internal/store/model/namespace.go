package model

import (
	"time"

	"github.com/odpf/siren/core/namespace"
)

type Namespace struct {
	ID          uint64 `gorm:"primarykey"`
	Provider    *Provider
	ProviderID  uint64 `gorm:"uniqueIndex:urn_provider_id_unique"`
	URN         string `gorm:"uniqueIndex:urn_provider_id_unique"`
	Name        string
	Credentials string
	Labels      StringStringMap `gorm:"type:jsonb" sql:"type:jsonb" `
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ns *Namespace) FromDomain(n *namespace.EncryptedNamespace) {
	ns.ID = n.ID
	ns.URN = n.URN
	ns.Name = n.Name
	ns.ProviderID = n.Provider
	ns.Credentials = n.Credentials
	ns.Labels = StringStringMap(n.Labels)
	ns.CreatedAt = n.CreatedAt
	ns.UpdatedAt = n.UpdatedAt
}

func (ns *Namespace) ToDomain() *namespace.EncryptedNamespace {
	return &namespace.EncryptedNamespace{
		Namespace: &namespace.Namespace{
			ID:        ns.ID,
			URN:       ns.URN,
			Name:      ns.Name,
			Provider:  ns.ProviderID,
			Labels:    ns.Labels,
			CreatedAt: ns.CreatedAt,
			UpdatedAt: ns.UpdatedAt,
		},
		Credentials: ns.Credentials,
	}
}
