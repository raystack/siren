package provider

import (
	"github.com/odpf/siren/domain"
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

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) ListProviders() ([]*domain.Provider, error) {
	providers, err := service.repository.List()
	if err != nil {
		return nil, err
	}

	domainProviders := make([]*domain.Provider, 0, len(providers))
	for i := 0; i < len(providers); i++ {
		provider, _ := providers[i].toDomain()
		domainProviders = append(domainProviders, provider)
	}
	return domainProviders, nil

}

func (service Service) CreateProvider(provider *domain.Provider) (*domain.Provider, error) {
	p := &Provider{}
	p, err := p.fromDomain(provider)
	if err != nil {
		return nil, err
	}

	newProvider, err := service.repository.Create(p)
	if err != nil {
		return nil, err
	}
	return newProvider.toDomain()
}

func (service Service) GetProvider(id uint64) (*domain.Provider, error) {
	provider, err := service.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return provider.toDomain()
}

func (service Service) UpdateProvider(provider *domain.Provider) (*domain.Provider, error) {
	w := &Provider{}
	w, err := w.fromDomain(provider)
	if err != nil {
		return nil, err
	}

	newProvider, err := service.repository.Update(w)
	if err != nil {
		return nil, err
	}
	return newProvider.toDomain()
}

func (service Service) DeleteProvider(id uint64) error {
	err := service.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
