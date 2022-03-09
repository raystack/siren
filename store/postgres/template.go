package postgres

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
	"gorm.io/gorm"
	"text/template"
)

const (
	leftDelim  = "[["
	rightDelim = "]]"
)

// TemplateRepository talks to the store to read or insert data
type TemplateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository returns repository struct
func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{db}
}

func (r TemplateRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Template{})
	if err != nil {
		return err
	}
	return nil
}

func (r TemplateRepository) Upsert(template *model.Template) (*model.Template, error) {
	var newTemplate, existingTemplate model.Template
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

func (r TemplateRepository) Index(tag string) ([]model.Template, error) {
	var templates []model.Template
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

func (r TemplateRepository) GetByName(name string) (*model.Template, error) {
	var template model.Template
	result := r.db.Where(fmt.Sprintf("name = '%s'", name)).Find(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &template, nil
}

func (r TemplateRepository) Delete(name string) error {
	var template model.Template
	result := r.db.Where("name = ?", name).Delete(&template)
	return result.Error
}

func enrichWithDefaults(variables []domain.Variable, requestVariables map[string]string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(variables); i++ {
		name := variables[i].Name
		defaultValue := variables[i].Default
		val, ok := requestVariables[name]
		if ok {
			result[name] = val
		} else {
			result[name] = defaultValue
		}
	}
	return result
}

var templateParser = template.New("parser").Delims(leftDelim, rightDelim).Parse

func (r TemplateRepository) Render(name string, requestVariables map[string]string) (string, error) {
	templateFromDB, err := r.GetByName(name)
	if err != nil {
		return "", err
	}
	if templateFromDB == nil {
		return "", errors.New("template not found")
	}
	convertedTemplate, err := templateFromDB.ToDomain()
	if err != nil {
		return "", err
	}
	enrichedVariables := enrichWithDefaults(convertedTemplate.Variables, requestVariables)
	var tpl bytes.Buffer
	tmpl, err := templateParser(convertedTemplate.Body)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&tpl, enrichedVariables)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
