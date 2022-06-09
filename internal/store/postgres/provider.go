package postgres

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"gorm.io/gorm"
)

type Filters struct {
	Urn  string `mapstructure:"urn" validate:"omitempty"`
	Type string `mapstructure:"type" validate:"omitempty"`
}

// ProviderRepository talks to the store to read or insert data
type ProviderRepository struct {
	db *gorm.DB
}

// NewProviderRepository returns repository struct
func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db}
}

func (r ProviderRepository) List(filters map[string]interface{}) ([]*provider.Provider, error) {
	var providers []*model.Provider
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
	domainProviders := make([]*provider.Provider, 0, len(providers))
	for i := 0; i < len(providers); i++ {
		prov := providers[i].ToDomain()
		domainProviders = append(domainProviders, prov)
	}
	return domainProviders, nil
}

func (r ProviderRepository) Create(provider *provider.Provider) (*provider.Provider, error) {
	var newProvider model.Provider
	result := r.db.Create(newProvider.FromDomain(provider))
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("id = %d", provider.Id)).Find(&newProvider)
	if result.Error != nil {
		return nil, result.Error
	}

	return newProvider.ToDomain(), nil
}

func (r ProviderRepository) Get(id uint64) (*provider.Provider, error) {
	var provider model.Provider
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&provider)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return provider.ToDomain(), nil
}

func (r ProviderRepository) Update(provider *provider.Provider) (*provider.Provider, error) {
	inputProvider := model.Provider{}
	inputProvider.FromDomain(provider)
	var newProvider, existingProvider model.Provider
	result := r.db.Where(fmt.Sprintf("id = %d", inputProvider.Id)).Find(&existingProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("provider doesn't exist")
	} else {
		result = r.db.Where("id = ?", inputProvider.Id).Updates(inputProvider)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", inputProvider.Id)).Find(&newProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	return newProvider.ToDomain(), nil
}

func (r ProviderRepository) Delete(id uint64) error {
	var provider model.Provider
	result := r.db.Where("id = ?", id).Delete(&provider)
	return result.Error
}

func (r ProviderRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Provider{})
	if err != nil {
		return err
	}
	return nil
}
