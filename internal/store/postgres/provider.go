package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
)

const providerInsertQuery = `
INSERT INTO providers (host, urn, name, type, credentials, labels, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, now(), now())
RETURNING *
`

const providerUpdateQuery = `
UPDATE providers SET host=$2, urn=$3, name=$4, type=$5, credentials=$6, labels=$7, updated_at=now()
WHERE id = $1
RETURNING *
`

var providerListQueryBuilder = sq.Select(
	"id",
	"host",
	"urn",
	"name",
	"type",
	"credentials",
	"labels",
	"created_at",
	"updated_at",
).From("providers")

const providerDeleteQuery = `
DELETE from providers where id=$1
`

// ProviderRepository talks to the store to read or insert data
type ProviderRepository struct {
	client *pgc.Client
}

// NewProviderRepository returns repository struct
func NewProviderRepository(client *pgc.Client) *ProviderRepository {
	return &ProviderRepository{client}
}

func (r ProviderRepository) List(ctx context.Context, flt provider.Filter) ([]provider.Provider, error) {
	var queryBuilder = providerListQueryBuilder
	if flt.URN != "" {
		queryBuilder = queryBuilder.Where("urn = ?", flt.URN)
	}

	if flt.Type != "" {
		queryBuilder = queryBuilder.Where("type = ?", flt.Type)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Select *"),
		attribute.String("db.sql.table", "providers"),
	}
	rows, err := r.client.QueryxContext(otelsql.AddMeterLabels(ctx, attrs...), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	providersDomain := []provider.Provider{}
	for rows.Next() {
		var providerModel model.Provider
		if err := rows.StructScan(&providerModel); err != nil {
			return nil, err
		}

		providersDomain = append(providersDomain, *providerModel.ToDomain())
	}

	return providersDomain, nil
}

func (r ProviderRepository) Create(ctx context.Context, prov *provider.Provider) error {
	if prov == nil {
		return errors.New("provider domain is nil")
	}

	var provModel model.Provider
	provModel.FromDomain(*prov)

	var createdProvider model.Provider
	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Insert"),
		attribute.String("db.sql.table", "providers"),
	}

	if err := r.client.QueryRowxContext(otelsql.AddMeterLabels(ctx, attrs...), providerInsertQuery,
		provModel.Host,
		provModel.URN,
		provModel.Name,
		provModel.Type,
		provModel.Credentials,
		provModel.Labels,
	).StructScan(&createdProvider); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, pgc.ErrDuplicateKey) {
			return provider.ErrDuplicate
		}
		return err
	}

	*prov = *createdProvider.ToDomain()

	return nil
}

func (r ProviderRepository) Get(ctx context.Context, id uint64) (*provider.Provider, error) {
	query, args, err := providerListQueryBuilder.
		Where("id = ?", id).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var provModel model.Provider

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Select"),
		attribute.String("db.sql.table", "providers"),
	}

	if err := r.client.GetContext(otelsql.AddMeterLabels(ctx, attrs...), &provModel, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, provider.NotFoundError{ID: id}
		}
		return nil, err
	}

	return provModel.ToDomain(), nil
}

func (r ProviderRepository) Update(ctx context.Context, provDomain *provider.Provider) error {
	if provDomain == nil {
		return errors.New("provider domain is nil")
	}

	var provModel model.Provider
	provModel.FromDomain(*provDomain)

	var updatedProvider model.Provider

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Update"),
		attribute.String("db.sql.table", "providers"),
	}

	if err := r.client.QueryRowxContext(otelsql.AddMeterLabels(ctx, attrs...), providerUpdateQuery,
		provModel.ID,
		provModel.Host,
		provModel.URN,
		provModel.Name,
		provModel.Type,
		provModel.Credentials,
		provModel.Labels,
	).StructScan(&updatedProvider); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return provider.NotFoundError{ID: provModel.ID}
		}
		if errors.Is(err, pgc.ErrDuplicateKey) {
			return provider.ErrDuplicate
		}
		return err
	}

	*provDomain = *updatedProvider.ToDomain()

	return nil
}

func (r ProviderRepository) Delete(ctx context.Context, id uint64) error {
	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Delete"),
		attribute.String("db.sql.table", "providers"),
	}
	if _, err := r.client.ExecContext(otelsql.AddMeterLabels(ctx, attrs...), providerDeleteQuery, id); err != nil {
		return err
	}
	return nil
}
