package provider

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service Service) ListProviders(filters map[string]interface{}) ([]*Provider, error) {
	return service.repository.List(filters)
}

func (service Service) CreateProvider(provider *Provider) (*Provider, error) {
	return service.repository.Create(provider)
}

func (service Service) GetProvider(id uint64) (*Provider, error) {
	return service.repository.Get(id)
}

func (service Service) UpdateProvider(provider *Provider) (*Provider, error) {
	return service.repository.Update(provider)
}

func (service Service) DeleteProvider(id uint64) error {
	return service.repository.Delete(id)
}
