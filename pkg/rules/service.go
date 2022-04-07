package rules

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/templates"
	"github.com/odpf/siren/store"
)

// Service handles business logic
type Service struct {
	repository      store.RuleRepository
	templateService domain.TemplatesService
}

// NewService returns repository struct
func NewService(repository store.RuleRepository, templateRepository store.TemplatesRepository) domain.RuleService {
	return &Service{
		repository:      repository,
		templateService: templates.NewService(templateRepository),
	}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Upsert(rule *domain.Rule) error {
	return service.repository.Upsert(rule, service.templateService)
}

func (service Service) Get(name, namespace, groupName, template string, providerNamespace uint64) ([]domain.Rule, error) {
	return service.repository.Get(name, namespace, groupName, template, providerNamespace)
}
