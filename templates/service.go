package templates

import (
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository *Repository
}

// NewService returns repository struct
func NewService(db *gorm.DB) *Service {
	return &Service{NewRepository(db)}
}

func (service *Service) Upsert(template *domain.Template) (*domain.Template, error) {
	t := &Template{}
	t, err := t.fromDomain(template)
	if err != nil {
		return nil, err
	}
	upsertedTemplate, err := service.repository.Upsert(t)
	if err != nil {
		return nil, err
	}
	return upsertedTemplate.toDomain()
}

func (service *Service) Index(tag string) ([]domain.Template, error) {
	templates, err := service.repository.Index(tag)
	if err != nil {
		return nil, err
	}
	domainTemplates := make([]domain.Template, 0, len(templates))
	for i := 0; i < len(templates); i++ {
		t, _ := templates[i].toDomain()
		domainTemplates = append(domainTemplates, *t)
	}
	return domainTemplates, nil
}

func (service *Service) GetByName(name string) (*domain.Template, error) {
	templates, err := service.repository.GetByName(name)
	if err != nil || templates == nil {
		return nil, err
	}
	return templates.toDomain()
}

func (service *Service) Delete(name string) error {
	return service.repository.Delete(name)
}

func (service *Service) Render(name string, body map[string]string) (string, error) {
	return service.repository.Render(name, body)
}
