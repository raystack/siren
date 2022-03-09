package rules

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
	"time"
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
	Variables             string           `gorm:"type:jsonb" sql:"type:jsonb"`
	ProviderNamespace     uint64           `gorm:"uniqueIndex:unique_name"`
	ProviderNamespaceInfo *model.Namespace `gorm:"foreignKey:ProviderNamespace"`
}

func (rule *Rule) fromDomain(r *domain.Rule) (*Rule, error) {
	rule.Id = r.Id
	rule.Name = r.Name
	rule.Enabled = &r.Enabled
	rule.GroupName = r.GroupName
	rule.Namespace = r.Namespace
	rule.Template = r.Template

	jsonString, err := json.Marshal(r.Variables)
	if err != nil {
		return nil, err
	}

	rule.Variables = string(jsonString)
	rule.ProviderNamespace = r.ProviderNamespace
	rule.CreatedAt = r.CreatedAt
	rule.UpdatedAt = r.UpdatedAt
	return rule, nil
}

func (rule *Rule) toDomain() (*domain.Rule, error) {
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

//Repository interface
type RuleRepository interface {
	Upsert(*Rule, domain.TemplatesService) (*Rule, error)
	Get(string, string, string, string, uint64) ([]Rule, error)
	Migrate() error
}
