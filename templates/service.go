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
