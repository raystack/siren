package template

import (
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname TemplateRepository --filename template_repository.go --output=./mocks
type Repository interface {
	Upsert(*Template) error
	Index(string) ([]Template, error)
	GetByName(string) (*Template, error)
	Delete(string) error
	Render(string, map[string]string) (string, error)
	Migrate() error
}

type Variable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Template struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Name      string     `json:"name" validate:"required"`
	Body      string     `json:"body" validate:"required"`
	Tags      []string   `json:"tags" validate:"required"`
	Variables []Variable `json:"variables" validate:"required,dive,required"`
}
