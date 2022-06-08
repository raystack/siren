package templates

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store"
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

func (service Service) Upsert(template *domain.Template) error {
	return service.repository.Upsert(template)
}

func (service Service) Index(tag string) ([]domain.Template, error) {
	return service.repository.Index(tag)
}

func (service Service) GetByName(name string) (*domain.Template, error) {
	return service.repository.GetByName(name)
}

func (service Service) Delete(name string) error {
	return service.repository.Delete(name)
}

func (service Service) Render(name string, body map[string]string) (string, error) {
	return service.repository.Render(name, body)
}
