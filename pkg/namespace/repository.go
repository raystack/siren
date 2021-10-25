package namespace

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) List() ([]*Namespace, error) {
	var namespaces []*Namespace
	selectQuery := fmt.Sprintf("select * from namespaces")
	result := r.db.Raw(selectQuery).Find(&namespaces)
	if result.Error != nil {
		return nil, result.Error
	}

	return namespaces, nil
}

func (r Repository) Create(namespace *Namespace) (*Namespace, error) {
	var newNamespace Namespace
	result := r.db.Create(namespace)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("id = %d", namespace.Id)).Find(&newNamespace)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newNamespace, nil
}

func (r Repository) Get(id uint64) (*Namespace, error) {
	var namespace Namespace
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&namespace)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &namespace, nil
}

func (r Repository) Update(namespace *Namespace) (*Namespace, error) {
	var newNamespace, existingNamespace Namespace
	result := r.db.Where(fmt.Sprintf("id = %d", namespace.Id)).Find(&existingNamespace)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("namespace doesn't exist")
	} else {
		result = r.db.Where("id = ?", namespace.Id).Updates(namespace)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", namespace.Id)).Find(&newNamespace)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newNamespace, nil
}

func (r Repository) Delete(id uint64) error {
	var namespace Namespace
	result := r.db.Where("id = ?", id).Delete(&namespace)
	return result.Error
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Namespace{})
	if err != nil {
		return err
	}
	return nil
}
