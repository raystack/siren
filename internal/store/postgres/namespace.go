package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

// NamespaceRepository talks to the store to read or insert data
type NamespaceRepository struct {
	client *Client
}

// NewNamespaceRepository returns repository struct
func NewNamespaceRepository(client *Client) *NamespaceRepository {
	return &NamespaceRepository{client}
}

func (r NamespaceRepository) List(ctx context.Context) ([]namespace.EncryptedNamespace, error) {
	var namespaceModels []*model.Namespace
	if err := r.client.db.WithContext(ctx).Raw("select * from namespaces").Find(&namespaceModels).Error; err != nil {
		return nil, err
	}

	var result []namespace.EncryptedNamespace
	for _, m := range namespaceModels {
		n, err := m.ToDomain()
		if err != nil {
			// log error here
			continue
		}
		result = append(result, *n)
	}
	return result, nil
}

func (r NamespaceRepository) Create(ctx context.Context, ns *namespace.EncryptedNamespace) (uint64, error) {
	nsModel := new(model.Namespace)
	if err := nsModel.FromDomain(ns); err != nil {
		return 0, err
	}

	if err := r.client.db.WithContext(ctx).Create(nsModel).Error; err != nil {
		err = checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return 0, namespace.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return 0, namespace.ErrRelation
		}
		return 0, err
	}

	return nsModel.ID, nil
}

func (r NamespaceRepository) Get(ctx context.Context, id uint64) (*namespace.EncryptedNamespace, error) {
	var nsModel model.Namespace
	result := r.client.db.WithContext(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&nsModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, namespace.NotFoundError{ID: id}
	}
	ns, err := nsModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (r NamespaceRepository) Update(ctx context.Context, ns *namespace.EncryptedNamespace) (uint64, error) {
	m := new(model.Namespace)
	if err := m.FromDomain(ns); err != nil {
		return 0, err
	}

	result := r.client.db.Where("id = ?", m.ID).Updates(m)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return 0, namespace.ErrDuplicate
		}
		if errors.Is(err, errForeignKeyViolation) {
			return 0, namespace.ErrRelation
		}
		return 0, err
	}
	if result.RowsAffected == 0 {
		return 0, namespace.NotFoundError{ID: ns.ID}
	}

	return m.ID, nil
}

func (r NamespaceRepository) Delete(ctx context.Context, id uint64) error {
	var namespace model.Namespace
	result := r.client.db.WithContext(ctx).Where("id = ?", id).Delete(&namespace)
	return result.Error
}
