package postgres

import (
	"context"

	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

type RuleResponse struct {
	NamespaceUrn string
	ProviderUrn  string
	ProviderType string
	ProviderHost string
}

// RuleRepository talks to the store to read or insert data
type RuleRepository struct {
	*transaction
}

// NewRuleRepository returns repository struct
func NewRuleRepository(db *gorm.DB) *RuleRepository {
	return &RuleRepository{&transaction{db}}
}

func (r *RuleRepository) Upsert(ctx context.Context, rl *rule.Rule) error {
	m := new(model.Rule)
	if err := m.FromDomain(rl); err != nil {
		return err
	}

	result := r.getDb(ctx).Where("name = ?", m.Name).Updates(m)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return rule.ErrDuplicate
		}
		return err
	}

	if result.RowsAffected == 0 {
		if err := r.getDb(ctx).Create(m).Error; err != nil {
			err = checkPostgresError(err)
			if errors.Is(err, errDuplicateKey) {
				return rule.ErrDuplicate
			}
			return err
		}
	}

	if err := r.getDb(ctx).Where("name = ?", m.Name).Find(m).Error; err != nil {
		return err
	}

	newRule, err := m.ToDomain()
	if err != nil {
		return err
	}
	*rl = *newRule
	return nil
}

func (r *RuleRepository) Get(ctx context.Context, name, namespace, groupName, template string, providerNamespace uint64) ([]rule.Rule, error) {
	var rules []model.Rule
	db := r.getDb(ctx)
	if name != "" {
		db = db.Where("name = ?", name)
	}
	if namespace != "" {
		db = db.Where("namespace = ?", namespace)
	}
	if groupName != "" {
		db = db.Where("group_name = ?", groupName)
	}
	if template != "" {
		db = db.Where("template = ?", template)
	}
	if providerNamespace != 0 {
		db = db.Where("provider_namespace = ?", providerNamespace)
	}

	if err := db.Find(&rules).Error; err != nil {
		return nil, err
	}

	var domainRules []rule.Rule
	for _, r := range rules {
		rule, err := r.ToDomain()
		if err != nil {
			return nil, err
		}
		domainRules = append(domainRules, *rule)
	}

	return domainRules, nil
}
