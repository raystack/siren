package rules

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"time"
)

type Rule struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"index:idx_rule_name,unique"`
	Namespace string
	Entity    string
	GroupName string
	Template  string
	Status    string
	Variables string `gorm:"type:jsonb" sql:"type:jsonb"`
}

func (rule *Rule) fromDomain(r *domain.Rule) (*Rule, error) {
	rule.ID = r.ID
	rule.CreatedAt = r.CreatedAt
	rule.UpdatedAt = r.UpdatedAt
	rule.Name = r.Name
	rule.Namespace = r.Namespace
	rule.GroupName = r.GroupName
	rule.Entity = r.Entity
	rule.Template = r.Template
	rule.Status = r.Status
	jsonString, err := json.Marshal(r.Variables)
	if err != nil {
		return nil, err
	}
	rule.Variables = string(jsonString)
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
		ID:        rule.ID,
		Name:      rule.Name,
		Namespace: rule.Namespace,
		Entity:    rule.Entity,
		GroupName: rule.GroupName,
		Template:  rule.Template,
		Status:    rule.Status,
		CreatedAt: rule.CreatedAt,
		UpdatedAt: rule.UpdatedAt,
		Variables: variables,
	}, nil
}

//Repository interface
type RuleRepository interface {
	Upsert(*Rule, cortexCaller, domain.TemplatesService) (*Rule, error)
	Get(string, string, string, string, string) ([]Rule, error)
	Migrate() error
}
