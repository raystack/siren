package provider

import "github.com/odpf/siren/domain"

// Service handling the business logics
type Service struct {
	providerRepository domain.ProviderRepository
}

// NewService returns service struct
func NewService(pr domain.ProviderRepository) *Service {
	return &Service{pr}
}

// Create record
func (s *Service) Create(p *domain.Provider) error {
	return s.providerRepository.Create(p)
}
