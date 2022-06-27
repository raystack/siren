package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

// TemplateRepository talks to the store to read or insert data
type TemplateRepository struct {
	client *Client
}

// NewTemplateRepository returns repository struct
func NewTemplateRepository(client *Client) *TemplateRepository {
	return &TemplateRepository{client}
}

func (r TemplateRepository) Upsert(ctx context.Context, tmpl *template.Template) (uint64, error) {
	modelTemplate := &model.Template{}
	err := modelTemplate.FromDomain(tmpl)
	if err != nil {
		return 0, err
	}

	result := r.client.db.WithContext(ctx).Where("name = ?", modelTemplate.Name).Updates(&modelTemplate)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return 0, template.ErrDuplicate
		}
		return 0, err
	}

	if result.RowsAffected == 0 {
		if err := r.client.db.WithContext(ctx).Create(&modelTemplate).Error; err != nil {
			err = checkPostgresError(err)
			if errors.Is(err, errDuplicateKey) {
				return 0, template.ErrDuplicate
			}
			return 0, err
		}
	}

	return modelTemplate.ID, err
}

func (r TemplateRepository) List(ctx context.Context, flt template.Filter) ([]template.Template, error) {
	var (
		templates []model.Template
		result    = r.client.db
	)
	if flt.Tag != "" {
		result = result.Where("tags @>ARRAY[?]", flt.Tag)
	}
	result = result.WithContext(ctx).Find(&templates)
	if result.Error != nil {
		return nil, result.Error
	}

	domainTemplates := make([]template.Template, 0, len(templates))
	for _, templateModel := range templates {
		templateDomain, err := templateModel.ToDomain()
		if err != nil {
			return nil, err
		}
		domainTemplates = append(domainTemplates, *templateDomain)
	}
	return domainTemplates, nil
}

func (r TemplateRepository) GetByName(ctx context.Context, name string) (*template.Template, error) {
	var templateModel model.Template
	result := r.client.db.WithContext(ctx).Where(fmt.Sprintf("name = '%s'", name)).Find(&templateModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, template.NotFoundError{Name: name}
	}
	tmpl, err := templateModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (r TemplateRepository) Delete(ctx context.Context, name string) error {
	var template model.Template
	result := r.client.db.WithContext(ctx).Where("name = ?", name).Delete(&template)
	return result.Error
}
