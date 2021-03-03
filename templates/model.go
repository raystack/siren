package templates

import (
	"github.com/lib/pq"
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Name string         `json:"name" gorm:"index:idx_name,unique"`
	Body string         `json:"body"`
	Tags pq.StringArray `gorm:"type:text[]" json:"tags"`
}

func (template *Template) fromDomain(t *domain.Template) (*Template, error) {
	template.ID = t.ID
	template.CreatedAt = t.CreatedAt
	template.UpdatedAt = t.UpdatedAt
	template.Name = t.Name
	template.Tags = t.Tags
	template.Body = t.Body
	return template, nil
}

func (template *Template) toDomain() (*domain.Template, error) {
	return &domain.Template{
		ID:        template.ID,
		Name:      template.Name,
		Body:      template.Body,
		Tags:      template.Tags,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
	}, nil
}
