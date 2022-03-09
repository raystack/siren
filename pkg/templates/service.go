package templates

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store"
	"github.com/odpf/siren/store/model"
)

// Service handles business logic
type Service struct {
	repository store.TemplatesRepository
}

// NewService returns repository struct
func NewService(repository store.TemplatesRepository) domain.TemplatesService {
	return &Service{repository}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Upsert(template *domain.Template) (*domain.Template, error) {
	t := &model.Template{}
	t, err := t.FromDomain(template)
	if err != nil {
		return nil, err
	}
	upsertedTemplate, err := service.repository.Upsert(t)
	if err != nil {
		return nil, err
	}
	return upsertedTemplate.ToDomain()
}

func (service Service) Index(tag string) ([]domain.Template, error) {
	templates, err := service.repository.Index(tag)
	if err != nil {
		return nil, err
	}
	domainTemplates := make([]domain.Template, 0, len(templates))
	for i := 0; i < len(templates); i++ {
		t, _ := templates[i].ToDomain()
		domainTemplates = append(domainTemplates, *t)
	}
	return domainTemplates, nil
}

func (service Service) GetByName(name string) (*domain.Template, error) {
	template, err := service.repository.GetByName(name)
	if err != nil || template == nil {
		return nil, err
	}
	return template.ToDomain()
}

func (service Service) Delete(name string) error {
	return service.repository.Delete(name)
}

func (service Service) Render(name string, body map[string]string) (string, error) {
	return service.repository.Render(name, body)
}
