package postgres

import (
	"fmt"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

// NamespaceRepository talks to the store to read or insert data
type NamespaceRepository struct {
	db *gorm.DB
}

// NewNamespaceRepository returns repository struct
func NewNamespaceRepository(db *gorm.DB) *NamespaceRepository {
	return &NamespaceRepository{db}
}

func (r NamespaceRepository) List() ([]*namespace.EncryptedNamespace, error) {
	var namespaceModels []*model.Namespace
	selectQuery := "select * from namespaces"
	if err := r.db.Raw(selectQuery).Find(&namespaceModels).Error; err != nil {
		return nil, err
	}

	var result []*namespace.EncryptedNamespace
	for _, m := range namespaceModels {
		n, err := m.ToDomain()
		if err != nil {
			// log error here
			continue
		}
		result = append(result, n)
	}
	return result, nil
}

func (r NamespaceRepository) Create(ns *namespace.EncryptedNamespace) error {
	nsModel := new(model.Namespace)
	if err := nsModel.FromDomain(ns); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(nsModel).Error; err != nil {
			err = checkPostgresError(err)
			if errors.Is(err, errDuplicateKey) {
				return namespace.ErrDuplicate
			}
			return err
		}

		newNamespace, err := nsModel.ToDomain()
		if err != nil {
			return err
		}
		*ns = *newNamespace
		return nil
	})
}

func (r NamespaceRepository) Get(id uint64) (*namespace.EncryptedNamespace, error) {
	var nsModel model.Namespace
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&nsModel)
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

func (r NamespaceRepository) Update(ns *namespace.EncryptedNamespace) error {
	m := new(model.Namespace)
	if err := m.FromDomain(ns); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		result := r.db.Where("id = ?", m.ID).Updates(m)
		if result.Error != nil {
			err := checkPostgresError(result.Error)
			if errors.Is(err, errDuplicateKey) {
				return namespace.ErrDuplicate
			}
			return err
		}
		if result.RowsAffected == 0 {
			return namespace.NotFoundError{ID: ns.ID}
		}
		//TODO need to check whether this is necessary or not
		if err := r.db.Where(fmt.Sprintf("id = %d", m.ID)).Find(m).Error; err != nil {
			return err
		}

		newNamespace, err := m.ToDomain()
		if err != nil {
			return err
		}
		*ns = *newNamespace
		return nil
	})
}

func (r NamespaceRepository) Delete(id uint64) error {
	var namespace model.Namespace
	result := r.db.Where("id = ?", id).Delete(&namespace)
	return result.Error
}
