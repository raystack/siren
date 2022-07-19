package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

// ProviderRepository talks to the store to read or insert data
type ProviderRepository struct {
	client *Client
}

// NewProviderRepository returns repository struct
func NewProviderRepository(client *Client) *ProviderRepository {
	return &ProviderRepository{client}
}

func (r ProviderRepository) List(ctx context.Context, flt provider.Filter) ([]provider.Provider, error) {
	var providers []*model.Provider

	db := r.client.db
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
		domainProviders = append(domainProviders, *provModel.ToDomain())
	}
	return domainProviders, nil
}

func (r ProviderRepository) Create(ctx context.Context, prov *provider.Provider) error {
	if prov == nil {
		return errors.New("provider domain is nil")
	}

	var provModel model.Provider
	provModel.FromDomain(prov)

	result := r.client.db.WithContext(ctx).Create(&provModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return provider.ErrDuplicate
		}
		return result.Error
	}

	return nil
}

func (r ProviderRepository) Get(ctx context.Context, id uint64) (*provider.Provider, error) {
	var provModel model.Provider
	result := r.client.db.WithContext(ctx).Where(fmt.Sprintf("id = %d", id)).Find(&provModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, provider.NotFoundError{ID: id}
	}

	return provModel.ToDomain(), nil
}

func (r ProviderRepository) Update(ctx context.Context, provDomain *provider.Provider) error {
	if provDomain == nil {
		return errors.New("provider domain is nil")
	}

	var provModel model.Provider
	provModel.FromDomain(provDomain)

	result := r.client.db.Where("id = ?", provModel.ID).Updates(&provModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errDuplicateKey) {
			return provider.ErrDuplicate
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return provider.NotFoundError{ID: provModel.ID}
	}

	return nil
}

func (r ProviderRepository) Delete(ctx context.Context, id uint64) error {
	var provider model.Provider
	result := r.client.db.WithContext(ctx).Where("id = ?", id).Delete(&provider)
	return result.Error
}
