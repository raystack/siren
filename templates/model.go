package templates

import (
	"encoding/json"
	"github.com/lib/pq"
	"github.com/odpf/siren/domain"
	"time"
)

type Template struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"index:idx_name,unique"`
	Body      string
	Tags      pq.StringArray `gorm:"type:text[]"`
	Variables string         `sql:"type:jsonb" gorm:"type:jsonb"`
}

func (template *Template) fromDomain(t *domain.Template) (*Template, error) {
	template.ID = t.ID
	template.CreatedAt = t.CreatedAt
	template.UpdatedAt = t.UpdatedAt
	template.Name = t.Name
	template.Tags = t.Tags
	template.Body = t.Body
	jsonString, err := json.Marshal(t.Variables)
	if err != nil {
		return nil, err
	}
	template.Variables = string(jsonString)
	return template, nil
}

func (template *Template) toDomain() (*domain.Template, error) {
	var variables []domain.Variable
	jsonBlob := []byte(template.Variables)
	err := json.Unmarshal(jsonBlob, &variables)
	if err != nil {
		return nil, err
	}
	return &domain.Template{
		ID:        template.ID,
		Name:      template.Name,
		Body:      template.Body,
		Tags:      template.Tags,
		CreatedAt: template.CreatedAt,
		UpdatedAt: template.UpdatedAt,
		Variables: variables,
	}, nil
}
