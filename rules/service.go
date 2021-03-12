package rules

import (
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository RuleRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.RuleService {
	return &Service{NewRepository(db)}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Upsert(rule *domain.Rule) (*domain.Rule, error) {
	r := &Rule{}
	r, err := r.fromDomain(rule)
	upsertedRule, err := service.repository.Upsert(r)
	if err != nil {
		return nil, err
	}
	return upsertedRule.toDomain()
}

func (service Service) Get(namespace, entity, groupName, status, template string) ([]domain.Rule, error) {
	rules, err := service.repository.Get(namespace, entity, groupName, status, template)
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
