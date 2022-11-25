package model

import (
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/pkg/pgc"
)

type Namespace struct {
	ID               uint64              `db:"id"`
	ProviderID       uint64              `db:"provider_id"`
	URN              string              `db:"urn"`
	Name             string              `db:"name"`
	CredentialString string              `db:"credentials"`
	Labels           pgc.StringStringMap `db:"labels"`
	CreatedAt        time.Time           `db:"created_at"`
	UpdatedAt        time.Time           `db:"updated_at"`
}

func (ns *Namespace) FromDomain(n namespace.EncryptedNamespace) {
	ns.ID = n.ID
	ns.URN = n.URN
	ns.Name = n.Name
	ns.ProviderID = n.Provider.ID
	ns.CredentialString = n.CredentialString
	ns.Labels = pgc.StringStringMap(n.Labels)
	ns.CreatedAt = n.CreatedAt
	ns.UpdatedAt = n.UpdatedAt
}

func (ns *Namespace) ToDomain() *namespace.EncryptedNamespace {
	return &namespace.EncryptedNamespace{
		Namespace: &namespace.Namespace{
			ID:   ns.ID,
			URN:  ns.URN,
			Name: ns.Name,
			Provider: provider.Provider{
				ID: ns.ProviderID,
			},
			Labels:    ns.Labels,
			CreatedAt: ns.CreatedAt,
			UpdatedAt: ns.UpdatedAt,
		},
		CredentialString: ns.CredentialString,
	}
}

type NamespaceDetail struct {
	ID               uint64              `db:"id"`
	Provider         Provider            `db:"provider"`
	URN              string              `db:"urn"`
	Name             string              `db:"name"`
	CredentialString string              `db:"credentials"`
	Labels           pgc.StringStringMap `db:"labels"`
	CreatedAt        time.Time           `db:"created_at"`
	UpdatedAt        time.Time           `db:"updated_at"`
}

func (ns *NamespaceDetail) FromDomain(n namespace.EncryptedNamespace) {
	ns.ID = n.ID
	ns.URN = n.URN
	ns.Name = n.Name
	ns.Provider.FromDomain(n.Provider)
	ns.CredentialString = n.CredentialString
	ns.Labels = pgc.StringStringMap(n.Labels)
	ns.CreatedAt = n.CreatedAt
	ns.UpdatedAt = n.UpdatedAt
}

func (ns *NamespaceDetail) ToDomain() *namespace.EncryptedNamespace {
	return &namespace.EncryptedNamespace{
		Namespace: &namespace.Namespace{
			ID:        ns.ID,
			URN:       ns.URN,
			Name:      ns.Name,
			Provider:  *ns.Provider.ToDomain(),
			Labels:    ns.Labels,
			CreatedAt: ns.CreatedAt,
			UpdatedAt: ns.UpdatedAt,
		},
		CredentialString: ns.CredentialString,
	}
}
