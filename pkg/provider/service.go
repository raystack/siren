package provider

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store"
)

// Service handles business logic
type Service struct {
	repository store.ProviderRepository
}

// NewService returns repository struct
func NewService(repository store.ProviderRepository) domain.ProviderService {
	return &Service{repository}
}

func (service Service) ListProviders(filters map[string]interface{}) ([]*domain.Provider, error) {
	return service.repository.List(filters)
}

func (service Service) CreateProvider(provider *domain.Provider) (*domain.Provider, error) {
	return service.repository.Create(provider)
}

func (service Service) GetProvider(id uint64) (*domain.Provider, error) {
	return service.repository.Get(id)
}

func (service Service) UpdateProvider(provider *domain.Provider) (*domain.Provider, error) {
	return service.repository.Update(provider)
}

func (service Service) DeleteProvider(id uint64) error {
	return service.repository.Delete(id)
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}
