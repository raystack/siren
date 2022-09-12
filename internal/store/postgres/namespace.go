package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const namespaceInsertQuery = `
INSERT INTO namespaces (provider_id, urn, name, credentials, labels, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *
`

const namespaceUpdateQuery = `
UPDATE namespaces SET provider_id=$2, urn=$3, name=$4, credentials=$5, labels=$6, updated_at=now()
WHERE id = $1
RETURNING *
`

var namespaceListQueryBuilder = sq.Select(
	"id",
	"provider_id",
	"urn",
	"name",
	"credentials",
	"labels",
	"created_at",
	"updated_at",
).From("namespaces")

const namespaceDeleteQuery = `
DELETE from namespaces where id=$1
`

// NamespaceRepository talks to the store to read or insert data
type NamespaceRepository struct {
	client *Client
}

// NewNamespaceRepository returns repository struct
func NewNamespaceRepository(client *Client) *NamespaceRepository {
	return &NamespaceRepository{client}
}

func (r NamespaceRepository) List(ctx context.Context) ([]namespace.EncryptedNamespace, error) {
	query, args, err := namespaceListQueryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.client.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var encryptedNamespaces []namespace.EncryptedNamespace
	for rows.Next() {
		var namespaceModel model.Namespace
		if err := rows.StructScan(&namespaceModel); err != nil {
			return nil, err
		}
		encryptedNamespaces = append(encryptedNamespaces, *namespaceModel.ToDomain())
	}

	return encryptedNamespaces, nil
}

func (r NamespaceRepository) Create(ctx context.Context, ns *namespace.EncryptedNamespace) error {
	if ns == nil {
		return errors.New("nil encrypted namespace domain when converting to namespace model")
	}

	nsModel := new(model.Namespace)
	nsModel.FromDomain(*ns)

	var createdNamespace model.Namespace
	if err := r.client.db.QueryRowxContext(ctx, namespaceInsertQuery,
		nsModel.ProviderID,
		nsModel.URN,
		nsModel.Name,
		nsModel.Credentials,
		nsModel.Labels,
	).StructScan(&createdNamespace); err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return namespace.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return namespace.ErrRelation
		}
		return err
	}

	*ns = *createdNamespace.ToDomain()

	return nil
}

func (r NamespaceRepository) Get(ctx context.Context, id uint64) (*namespace.EncryptedNamespace, error) {
	query, args, err := namespaceListQueryBuilder.Where("id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var nsModel model.Namespace
	if err := r.client.db.GetContext(ctx, &nsModel, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, namespace.NotFoundError{ID: id}
		}
		return nil, err
	}

	return nsModel.ToDomain(), nil
}

func (r NamespaceRepository) Update(ctx context.Context, ns *namespace.EncryptedNamespace) error {
	if ns == nil {
		return errors.New("nil encrypted namespace domain when converting to namespace model")
	}

	namespaceModel := new(model.Namespace)
	namespaceModel.FromDomain(*ns)

	var updatedNamespace model.Namespace
	if err := r.client.db.QueryRowxContext(ctx, namespaceUpdateQuery,
		namespaceModel.ID,
		namespaceModel.ProviderID,
		namespaceModel.URN,
		namespaceModel.Name,
		namespaceModel.Credentials,
		namespaceModel.Labels,
	).StructScan(&updatedNamespace); err != nil {
		err := checkPostgresError(err)
		if errors.Is(err, sql.ErrNoRows) {
			return namespace.NotFoundError{ID: namespaceModel.ID}
		}
		if errors.Is(err, errDuplicateKey) {
			return namespace.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return namespace.ErrRelation
		}
		return err
	}

	*ns = *updatedNamespace.ToDomain()

	return nil
}

func (r NamespaceRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.client.db.ExecContext(ctx, namespaceDeleteQuery, id); err != nil {
		return err
	}
	return nil
}
