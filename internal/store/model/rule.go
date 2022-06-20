package model

import (
	"encoding/json"
	"time"

	"github.com/odpf/siren/core/rule"
)

type Rule struct {
	Id                    uint64 `gorm:"primarykey"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Name                  string `gorm:"index:idx_rule_name,unique"`
	Namespace             string `gorm:"uniqueIndex:unique_name"`
	GroupName             string `gorm:"uniqueIndex:unique_name"`
	Template              string `gorm:"uniqueIndex:unique_name"`
	Enabled               *bool
	Variables             string     `gorm:"type:jsonb" sql:"type:jsonb"`
	ProviderNamespace     uint64     `gorm:"uniqueIndex:unique_name"`
	ProviderNamespaceInfo *Namespace `gorm:"foreignKey:ProviderNamespace"`
}

func (rl *Rule) FromDomain(r *rule.Rule) error {
	rl.Id = r.Id
	rl.Name = r.Name
	rl.Enabled = &r.Enabled
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
	var variables []rule.RuleVariable
	jsonBlob := []byte(rl.Variables)
	err := json.Unmarshal(jsonBlob, &variables)
	if err != nil {
		return nil, err
	}
	return &rule.Rule{
		Id:                rl.Id,
		Name:              rl.Name,
		Enabled:           *rl.Enabled,
		GroupName:         rl.GroupName,
		Namespace:         rl.Namespace,
		Template:          rl.Template,
		Variables:         variables,
		ProviderNamespace: rl.ProviderNamespace,
		CreatedAt:         rl.CreatedAt,
		UpdatedAt:         rl.UpdatedAt,
	}, nil
}
