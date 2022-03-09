package rules

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/templates"
	"github.com/odpf/siren/store"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository      RuleRepository
	templateService domain.TemplatesService
}

// NewService returns repository struct
func NewService(templateRepository store.TemplatesRepository, db *gorm.DB) domain.RuleService {
	return &Service{
		repository:      NewRepository(db),
		templateService: templates.NewService(templateRepository),
	}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Upsert(rule *domain.Rule) (*domain.Rule, error) {
	r := &Rule{}
	r, err := r.fromDomain(rule)
	if err != nil {
		return nil, err
	}
	upsertedRule, err := service.repository.Upsert(r, service.templateService)
	if err != nil {
		return nil, err
	}
	return upsertedRule.toDomain()
}

func (service Service) Get(name, namespace, groupName, template string, providerNamespace uint64) ([]domain.Rule, error) {
	rules, err := service.repository.Get(name, namespace, groupName, template, providerNamespace)
	if err != nil {
		return nil, err
	}
	domainRules := make([]domain.Rule, 0, len(rules))
	for i := 0; i < len(rules); i++ {
		r, _ := rules[i].toDomain()
		domainRules = append(domainRules, *r)
	}
	return domainRules, nil
}
