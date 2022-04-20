package postgres

import (
	"context"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
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

func (r *RuleRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Rule{})
	if err != nil {
		return err
	}
	return nil
}

func (r *RuleRepository) Upsert(ctx context.Context, rule *domain.Rule) error {
	m := new(model.Rule)
	if err := m.FromDomain(rule); err != nil {
		return err
	}

	if result := r.getDb(ctx).Where("name = ?", m.Name).Updates(m); result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		if err := r.getDb(ctx).Create(m).Error; err != nil {
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
	*rule = *newRule
	return nil
}

func (r *RuleRepository) Get(ctx context.Context, name, namespace, groupName, template string, providerNamespace uint64) ([]domain.Rule, error) {
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

	var domainRules []domain.Rule
	for _, r := range rules {
		rule, err := r.ToDomain()
		if err != nil {
			return nil, err
		}
		domainRules = append(domainRules, *rule)
	}

	return domainRules, nil
}
