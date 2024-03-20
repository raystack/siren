package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
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
	client *pgc.Client
}

// NewTemplateRepository returns repository struct
func NewTemplateRepository(client *pgc.Client) *TemplateRepository {
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

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Insert"),
		attribute.String("db.sql.table", "templates"),
	}
	
	if err := r.client.QueryRowxContext(otelsql.AddMeterLabels(ctx, attrs...), templateUpsertQuery,
		templateModel.Name,
		templateModel.Body,
		templateModel.Tags,
		templateModel.Variables,
	).StructScan(&upsertedTemplate); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, pgc.ErrDuplicateKey) {
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

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Select *"),
		attribute.String("db.sql.table", "templates"),
	}

	rows, err := r.client.QueryxContext(otelsql.AddMeterLabels(ctx, attrs...), query, args...)
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

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Select"),
		attribute.String("db.sql.table", "templates"),
	}

	if err = r.client.GetContext(otelsql.AddMeterLabels(ctx, attrs...), &templateModel, query, args...); err != nil {
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
	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Delete"),
		attribute.String("db.sql.table", "templates"),
	}

	if _, err := r.client.ExecContext(otelsql.AddMeterLabels(ctx, attrs...), templateDeleteByNameQuery, name); err != nil {
		return err
	}
	return nil
}
