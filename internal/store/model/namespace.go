package model

import (
	"time"

	"github.com/odpf/siren/core/namespace"
)

type Namespace struct {
	ID          uint64          `db:"id"`
	ProviderID  uint64          `db:"provider_id"`
	URN         string          `db:"urn"`
	Name        string          `db:"name"`
	Credentials string          `db:"credentials"`
	Labels      StringStringMap `db:"labels"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

func (ns *Namespace) FromDomain(n namespace.EncryptedNamespace) {
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
