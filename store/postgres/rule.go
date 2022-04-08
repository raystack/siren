package postgres

import (
	"context"
	"fmt"

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

func (r *RuleRepository) Upsert(ctx context.Context, rule *domain.Rule, templatesService domain.TemplatesService) error {
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
	selectQuery := `SELECT * from rules`
	selectQueryWithWhereClause := `SELECT * from rules WHERE `
	var filterConditions []string
	if name != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("name = '%s' ", name))
	}
	if namespace != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("namespace = '%s' ", namespace))
	}
	if groupName != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("group_name = '%s' ", groupName))
	}
	if template != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("template = '%s' ", template))
	}
	if providerNamespace != 0 {
		filterConditions = append(filterConditions, fmt.Sprintf("provider_namespace = '%d' ", providerNamespace))
	}
	var finalSelectQuery string
	if len(filterConditions) == 0 {
		finalSelectQuery = selectQuery
	} else {
		finalSelectQuery = selectQueryWithWhereClause
		for i := 0; i < len(filterConditions); i++ {
			if i == 0 {
				finalSelectQuery += filterConditions[i]
			} else {
				finalSelectQuery += " AND " + filterConditions[i]
			}
		}
	}
	result := r.getDb(ctx).Raw(finalSelectQuery).Scan(&rules)
	if result.Error != nil {
		return nil, result.Error
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

func (r *RuleRepository) ListByGroup(ctx context.Context, namespace, groupName string, providerNamespace uint64) ([]*domain.Rule, error) {
	var models []*model.Rule
	if err := r.getDb(ctx).Where(fmt.Sprintf("namespace = '%s' AND group_name = '%s' AND provider_namespace = '%d'",
		namespace, groupName, providerNamespace)).Find(&models).Error; err != nil {
		return nil, err
	}

	var rules []*domain.Rule
	for _, r := range models {
		rule, err := r.ToDomain()
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
