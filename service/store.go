package service

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store"
)

type TemplatesStore interface {
	Upsert(template *store.Template) (*domain.Template, error)
}
