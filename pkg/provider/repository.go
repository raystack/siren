package provider

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type Filters struct {
	Urn  string `mapstructure:"urn" validate:"omitempty"`
	Type string `mapstructure:"type" validate:"omitempty"`
}

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) List(filters map[string]interface{}) ([]*Provider, error) {
	var providers []*Provider
	var conditions Filters
	if err := mapstructure.Decode(filters, &conditions); err != nil {
		return nil, err
	}

	db := r.db
	if conditions.Urn != "" {
		db = db.Where(`"urn" = ?`, conditions.Urn)
	}
	if conditions.Type != "" {
		db = db.Where(`"type" = ?`, conditions.Type)
	}

	result := db.Find(&providers)
	if result.Error != nil {
		return nil, result.Error
	}

	return providers, nil
}

func (r Repository) Create(provider *Provider) (*Provider, error) {
	var newProvider Provider
	result := r.db.Create(provider)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("id = %d", provider.Id)).Find(&newProvider)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newProvider, nil
}

func (r Repository) Get(id uint64) (*Provider, error) {
	var provider Provider
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &provider, nil
}

func (r Repository) Update(provider *Provider) (*Provider, error) {
	var newProvider, existingProvider Provider
	result := r.db.Where(fmt.Sprintf("id = %d", provider.Id)).Find(&existingProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("provider doesn't exist")
	} else {
		result = r.db.Where("id = ?", provider.Id).Updates(provider)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", provider.Id)).Find(&newProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newProvider, nil
}

func (r Repository) Delete(id uint64) error {
	var provider Provider
	result := r.db.Where("id = ?", id).Delete(&provider)
	return result.Error
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Provider{})
	if err != nil {
		return err
	}
	return nil
}
