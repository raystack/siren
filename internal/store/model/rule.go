package model

import (
	"encoding/json"
	"time"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/pkg/errors"
)

type Rule struct {
	ID                uint64    `db:"id"`
	Name              string    `db:"name"`
	Namespace         string    `db:"namespace"`
	GroupName         string    `db:"group_name"`
	Template          string    `db:"template"`
	Enabled           bool      `db:"enabled"`
	Variables         string    `db:"variables"`
	ProviderNamespace uint64    `db:"provider_namespace"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (rl *Rule) FromDomain(r rule.Rule) error {
	rl.ID = r.ID
	rl.Name = r.Name
	rl.Enabled = r.Enabled
	rl.GroupName = r.GroupName
	rl.Namespace = r.Namespace
	rl.Template = r.Template

	jsonString, err := json.Marshal(r.Variables)
	if err != nil {
		return err
	}

	rl.Variables = string(jsonString)
	rl.ProviderNamespace = r.ProviderNamespace
	rl.CreatedAt = r.CreatedAt
	rl.UpdatedAt = r.UpdatedAt
	return nil
}

func (rl *Rule) ToDomain() (*rule.Rule, error) {
	if rl == nil {
		return nil, errors.New("rule model is nil")
	}
	var variables []rule.RuleVariable
	jsonBlob := []byte(rl.Variables)
	err := json.Unmarshal(jsonBlob, &variables)
	if err != nil {
		return nil, err
	}
	return &rule.Rule{
		ID:                rl.ID,
		Name:              rl.Name,
		Enabled:           rl.Enabled,
		GroupName:         rl.GroupName,
		Namespace:         rl.Namespace,
		Template:          rl.Template,
		Variables:         variables,
		ProviderNamespace: rl.ProviderNamespace,
		CreatedAt:         rl.CreatedAt,
		UpdatedAt:         rl.UpdatedAt,
	}, nil
}
