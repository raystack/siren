package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const templateUpsertQuery = `
INSERT INTO templates (name, body, tags, variables, created_at, updated_at)
	VALUES ($1, $2, $3, $4, now(), now())
ON CONFLICT (name) 
DO
	UPDATE SET body=$2, tags=$3, variables=$4, updated_at=now()
RETURNING *
`

const templateDeleteByNameQuery = `
DELETE from templates where name=$1
`

var templateListQueryBuilder = sq.Select(
	"id",
	"name",
	"body",
	"tags",
	"variables",
	"created_at",
	"updated_at",
).From("templates")

// TemplateRepository talks to the store to read or insert data
type TemplateRepository struct {
	client *Client
}

// NewTemplateRepository returns repository struct
func NewTemplateRepository(client *Client) *TemplateRepository {
	return &TemplateRepository{client}
}

func (r TemplateRepository) Upsert(ctx context.Context, tmpl *template.Template) error {
	if tmpl == nil {
		return errors.New("template domain is nil")
	}

	templateModel := new(model.Template)
	if err := templateModel.FromDomain(*tmpl); err != nil {
		return err
	}

	var upsertedTemplate model.Template
	if err := r.client.db.QueryRowxContext(ctx, templateUpsertQuery,
		templateModel.Name,
		templateModel.Body,
		templateModel.Tags,
		templateModel.Variables,
	).StructScan(&upsertedTemplate); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return template.ErrDuplicate
		}
		return err
	}

	newTemplate, err := upsertedTemplate.ToDomain()
	if err != nil {
		return err
	}

	*tmpl = *newTemplate

	return nil
}

func (r TemplateRepository) List(ctx context.Context, flt template.Filter) ([]template.Template, error) {
	var queryBuilder = templateListQueryBuilder
	if flt.Tag != "" {
		queryBuilder = queryBuilder.Where("tags @>ARRAY[?]", flt.Tag)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.client.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templatesDomain := []template.Template{}
	for rows.Next() {
		var templateModel model.Template
		if err := rows.StructScan(&templateModel); err != nil {
			return nil, err
		}
		td, err := templateModel.ToDomain()
		if err != nil {
			return nil, err
		}
		templatesDomain = append(templatesDomain, *td)
	}

	return templatesDomain, nil
}

func (r TemplateRepository) GetByName(ctx context.Context, name string) (*template.Template, error) {
	query, args, err := templateListQueryBuilder.Where("name = ?", name).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var templateModel model.Template
	if err := r.client.db.GetContext(ctx, &templateModel, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, template.NotFoundError{Name: name}
		}
		return nil, err
	}

	tmpl, err := templateModel.ToDomain()
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (r TemplateRepository) Delete(ctx context.Context, name string) error {
	if _, err := r.client.db.ExecContext(ctx, templateDeleteByNameQuery, name); err != nil {
		return err
	}
	return nil
}
