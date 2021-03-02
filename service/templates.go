package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store"
)

type TemplatesService interface {
	Upsert(*store.Template) (*domain.Template, error)
}

type templatesService struct {
	templatesStore TemplatesStore
}

func NewTemplatesService(templatesStore TemplatesStore) *templatesService {
	return &templatesService{
		templatesStore: templatesStore,
	}
}

//Upsert upserts templates based on template name
func (s *templatesService) Upsert(template *store.Template) (*domain.Template, error) {
	return s.templatesStore.Upsert(template)
}
