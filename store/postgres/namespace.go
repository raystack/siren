package postgres

import (
	"errors"
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
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

func (r NamespaceRepository) List() ([]*domain.EncryptedNamespace, error) {
	var namespaceModels []*model.Namespace
	selectQuery := "select * from namespaces"
	if err := r.db.Raw(selectQuery).Find(&namespaceModels).Error; err != nil {
		return nil, err
	}

	var result []*domain.EncryptedNamespace
	for _, m := range namespaceModels {
		n, err := m.ToDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func (r NamespaceRepository) Create(namespace *domain.EncryptedNamespace) error {
	nsModel := new(model.Namespace)
	if err := nsModel.FromDomain(namespace); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.Create(nsModel).Error; err != nil {
			return err
		}

		newNamespace, err := nsModel.ToDomain()
		if err != nil {
			return err
		}
		*namespace = *newNamespace

		return nil
	})
}

func (r NamespaceRepository) Get(id uint64) (*domain.EncryptedNamespace, error) {
	var namespace model.Namespace
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&namespace)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return namespace.ToDomain()
}

func (r NamespaceRepository) Update(namespace *domain.EncryptedNamespace) error {
	m := new(model.Namespace)
	if err := m.FromDomain(namespace); err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		result := r.db.Where("id = ?", m.Id).Updates(m)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("namespace doesn't exist")
		}

		newNamespace, err := m.ToDomain()
		if err != nil {
			return err
		}
		*namespace = *newNamespace

		return nil
	})
}

func (r NamespaceRepository) Delete(id uint64) error {
	var namespace model.Namespace
	result := r.db.Where("id = ?", id).Delete(&namespace)
	return result.Error
}

func (r NamespaceRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Namespace{})
	if err != nil {
		return err
	}
	return nil
}
