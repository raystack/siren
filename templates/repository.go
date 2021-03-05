package templates

import (
	"bytes"
	"fmt"
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
	"text/template"
)

const (
	leftDelim  = "[["
	rightDelim = "]]"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Template{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Upsert(template *Template) (*Template, error) {
	var newTemplate, existingTemplate Template
	result := r.db.Where(fmt.Sprintf("name = '%s'", template.Name)).Find(&existingTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		result = r.db.Create(template)
	} else {
		result = r.db.Where("id = ?", existingTemplate.ID).Updates(template)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	result = r.db.Where(fmt.Sprintf("name = '%s'", template.Name)).Find(&newTemplate)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newTemplate, nil
}

func (r Repository) Index(tag string) ([]Template, error) {
	var templates []Template
	var result *gorm.DB
	if tag == "" {
		result = r.db.Find(&templates)
	} else {
		result = r.db.Where("tags @>ARRAY[?]", tag).Find(&templates)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return templates, nil
}

func (r Repository) GetByName(name string) (*Template, error) {
	var template Template
	result := r.db.Where(fmt.Sprintf("name = '%s'", name)).Find(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &template, nil
}

func (r Repository) Delete(name string) error {
	var template Template
	result := r.db.Where("name = ?", name).Delete(&template)
	return result.Error
}

func enrichWithDefaults(variables []domain.Variable, body map[string]string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(variables); i++ {
		name := variables[i].Name
		defaultValue := variables[i].Default
		val, ok := body[name]
		if ok {
			result[name] = val
		} else {
			result[name] = defaultValue
		}
	}
	return result
}

var templateParser = template.New("test").Delims(leftDelim, rightDelim).Parse

func (r Repository) Render(name string, body map[string]string) (string, error) {
	templateFromDB, err := r.GetByName(name)
	if err != nil {
		return "", err
	}
	convertedTemplate, err := templateFromDB.toDomain()
	enrichedBody := enrichWithDefaults(convertedTemplate.Variables, body)
	var tpl bytes.Buffer
	tmpl, err := templateParser(convertedTemplate.Body)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&tpl, enrichedBody)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
