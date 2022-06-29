package provider

import "github.com/odpf/siren/pkg/errors"

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
	//TODO check provider is nil
	prov, err := service.repository.Create(provider)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return nil, errors.ErrConflict.WithMsgf(err.Error())
		}
		return nil, err
	}
	return prov, nil
}

func (service Service) GetProvider(id uint64) (*Provider, error) {
	prov, err := service.repository.Get(id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	return prov, nil
}

func (service Service) UpdateProvider(provider *Provider) (*Provider, error) {
	prov, err := service.repository.Update(provider)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return nil, errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	return prov, nil
}

func (service Service) DeleteProvider(id uint64) error {
	return service.repository.Delete(id)
}
