package provider

import (
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Service handles business logic
type Service struct {
	repository ProviderRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.ProviderService {
	return &Service{NewRepository(db)}
}

func (service Service) ListProviders(filters map[string]interface{}) ([]*domain.Provider, error) {
	providers, err := service.repository.List(filters)
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.List")
	}

	domainProviders := make([]*domain.Provider, 0, len(providers))
	for i := 0; i < len(providers); i++ {
		provider := providers[i].toDomain()
		domainProviders = append(domainProviders, provider)
	}

	return domainProviders, nil

}

func (service Service) CreateProvider(provider *domain.Provider) (*domain.Provider, error) {
	p := &Provider{}
	newProvider, err := service.repository.Create(p.fromDomain(provider))
	if err != nil {
		return nil, errors.Wrap(err, "service.repository.Create")
	}

	return newProvider.toDomain(), nil
}

func (service Service) GetProvider(id uint64) (*domain.Provider, error) {
	provider, err := service.repository.Get(id)
	if err != nil {
		return nil, err
	}

	return provider.toDomain(), nil
}

func (service Service) UpdateProvider(provider *domain.Provider) (*domain.Provider, error) {
	w := &Provider{}
	newProvider, err := service.repository.Update(w.fromDomain(provider))
	if err != nil {
		return nil, err
	}

	return newProvider.toDomain(), nil
}

func (service Service) DeleteProvider(id uint64) error {
	return service.repository.Delete(id)
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}
