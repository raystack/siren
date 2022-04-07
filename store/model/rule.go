package model

import (
	"encoding/json"
	"time"

	"github.com/odpf/siren/domain"
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

func (rule *Rule) FromDomain(r *domain.Rule) error {
	rule.Id = r.Id
	rule.Name = r.Name
	rule.Enabled = &r.Enabled
	rule.GroupName = r.GroupName
	rule.Namespace = r.Namespace
	rule.Template = r.Template

	jsonString, err := json.Marshal(r.Variables)
	if err != nil {
		return err
	}

	rule.Variables = string(jsonString)
	rule.ProviderNamespace = r.ProviderNamespace
	rule.CreatedAt = r.CreatedAt
	rule.UpdatedAt = r.UpdatedAt
	return nil
}

func (rule *Rule) ToDomain() (*domain.Rule, error) {
	var variables []domain.RuleVariable
	jsonBlob := []byte(rule.Variables)
	err := json.Unmarshal(jsonBlob, &variables)
	if err != nil {
		return nil, err
	}
	return &domain.Rule{
		Id:                rule.Id,
		Name:              rule.Name,
		Enabled:           *rule.Enabled,
		GroupName:         rule.GroupName,
		Namespace:         rule.Namespace,
		Template:          rule.Template,
		Variables:         variables,
		ProviderNamespace: rule.ProviderNamespace,
		CreatedAt:         rule.CreatedAt,
		UpdatedAt:         rule.UpdatedAt,
	}, nil
}
