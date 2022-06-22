package postgres

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

type Filters struct {
	URN  string `mapstructure:"urn" validate:"omitempty"`
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
	if conditions.URN != "" {
		db = db.Where(`"urn" = ?`, conditions.URN)
	}
	if conditions.Type != "" {
		db = db.Where(`"type" = ?`, conditions.Type)
	}

	result := db.Find(&providers)
	if result.Error != nil {
		return nil, result.Error
	}
	domainProviders := make([]*provider.Provider, 0, len(providers))
	for _, provModel := range providers {
		provDomain, err := provModel.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		domainProviders = append(domainProviders, provDomain)
	}
	return domainProviders, nil
}

func (r ProviderRepository) Create(prov *provider.Provider) (*provider.Provider, error) {
	var provModel model.Provider
	if err := provModel.FromDomain(prov); err != nil {
		return nil, err
	}
	result := r.db.Create(&provModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return nil, provider.ErrDuplicate
		}
		return nil, result.Error
	}

	//TODO need to check whether this is necessary or not
	result = r.db.Where(fmt.Sprintf("id = %d", prov.ID)).Find(&provModel)
	if result.Error != nil {
		return nil, result.Error
	}

	newProvider, err := provModel.ToDomain()
	if err != nil {
		return nil, errors.New("failed to convert provider from model to domain")
	}
	return newProvider, nil
}

func (r ProviderRepository) Get(id uint64) (*provider.Provider, error) {
	var provModel model.Provider
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&provModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, provider.NotFoundError{ID: id}
	}
	provDomain, err := provModel.ToDomain()
	if err != nil {
		return nil, errors.New("failed to convert provider from model to domain")
	}
	return provDomain, nil
}

func (r ProviderRepository) Update(provDomain *provider.Provider) (*provider.Provider, error) {
	var provModel model.Provider
	if err := provModel.FromDomain(provDomain); err != nil {
		return nil, err
	}
	var newProvider, existingProvider model.Provider
	result := r.db.Where(fmt.Sprintf("id = %d", provModel.ID)).Find(&existingProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, provider.NotFoundError{ID: provDomain.ID}
	} else {
		result = r.db.Where("id = ?", provModel.ID).Updates(provModel)
		if result.Error != nil {
			err := checkPostgresError(result.Error)
			if errors.Is(err, errDuplicateKey) {
				return nil, provider.ErrDuplicate
			}
			return nil, result.Error
		}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", provModel.ID)).Find(&newProvider)
	if result.Error != nil {
		return nil, result.Error
	}
	provDomain, err := provModel.ToDomain()
	if err != nil {
		return nil, errors.New("failed to convert provider from model to domain")
	}
	return provDomain, nil
}

func (r ProviderRepository) Delete(id uint64) error {
	var provider model.Provider
	result := r.db.Where("id = ?", id).Delete(&provider)
	return result.Error
}
