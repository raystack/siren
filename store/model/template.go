package model

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
	Name      string `gorm:"index:idx_template_name,unique"`
	Body      string
	Tags      pq.StringArray `gorm:"type:text[];index:idx_tags,type:gin"`
	Variables string         `gorm:"type:jsonb" sql:"type:jsonb" `
}

func (template *Template) FromDomain(t *domain.Template) error {
	template.ID = t.ID
	template.CreatedAt = t.CreatedAt
	template.UpdatedAt = t.UpdatedAt
	template.Name = t.Name
	template.Tags = t.Tags
	template.Body = t.Body
	jsonString, err := json.Marshal(t.Variables)
	if err != nil {
		return err
	}
	template.Variables = string(jsonString)
	return nil
}

func (template *Template) ToDomain() (*domain.Template, error) {
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
