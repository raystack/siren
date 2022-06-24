package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"gorm.io/gorm"
)

// ProviderRepository talks to the store to read or insert data
type ProviderRepository struct {
	db *gorm.DB
}

// NewProviderRepository returns repository struct
func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db}
}

func (r ProviderRepository) List(ctx context.Context, flt provider.Filter) ([]provider.Provider, error) {
	var providers []*model.Provider

	db := r.db
	if flt.URN != "" {
		db = db.Where(`"urn" = ?`, flt.URN)
	}
	if flt.Type != "" {
		db = db.Where(`"type" = ?`, flt.Type)
	}

	result := db.WithContext(ctx).Find(&providers)
	if result.Error != nil {
		return nil, result.Error
	}
	domainProviders := make([]provider.Provider, 0, len(providers))
	for _, provModel := range providers {
		provDomain, err := provModel.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		domainProviders = append(domainProviders, *provDomain)
	}
	return domainProviders, nil
}

func (r ProviderRepository) Create(ctx context.Context, prov *provider.Provider) (uint64, error) {
	var provModel model.Provider
	if err := provModel.FromDomain(prov); err != nil {
		return 0, err
	}
	result := r.db.WithContext(ctx).Create(&provModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return 0, provider.ErrDuplicate
		}
		return 0, result.Error
	}

	return provModel.ID, nil
}

func (r ProviderRepository) Get(ctx context.Context, id uint64) (*provider.Provider, error) {
	var provModel model.Provider
	result := r.db.WithContext(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&provModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, provider.NotFoundError{ID: id}
	}
	provDomain, err := provModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return provDomain, nil
}

func (r ProviderRepository) Update(ctx context.Context, provDomain *provider.Provider) (uint64, error) {
	var provModel model.Provider
	if err := provModel.FromDomain(provDomain); err != nil {
		return 0, err
	}

	result := r.db.Where("id = ?", provModel.ID).Updates(&provModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return 0, provider.ErrDuplicate
		}
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, provider.NotFoundError{ID: provModel.ID}
	}

	return provModel.ID, nil
}

func (r ProviderRepository) Delete(ctx context.Context, id uint64) error {
	var provider model.Provider
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&provider)
	return result.Error
}
