package template

import (
	"context"
	"time"
)

const (
	ReservedName_SystemDefault = "system-default"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname TemplateRepository --filename template_repository.go --output=./mocks
type Repository interface {
	Upsert(context.Context, *Template) error
	List(context.Context, Filter) ([]Template, error)
	GetByName(context.Context, string) (*Template, error)
	Delete(context.Context, string) error
}

type Variable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Template struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name" validate:"required"`
	Body      string     `json:"body" validate:"required"`
	Tags      []string   `json:"tags" validate:"required"`
	Variables []Variable `json:"variables" validate:"required,dive,required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func IsReservedName(templateName string) bool {
	return (templateName == ReservedName_SystemDefault)
}
