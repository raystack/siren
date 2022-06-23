package provider

import (
	"context"

	"github.com/odpf/siren/pkg/errors"
)

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (s *Service) List(ctx context.Context, flt Filter) ([]*Provider, error) {
	return s.repository.List(ctx, flt)
}

func (s *Service) Create(ctx context.Context, prov *Provider) (uint64, error) {
	//TODO check provider is nil
	id, err := s.repository.Create(ctx, prov)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return 0, errors.ErrConflict.WithMsgf(err.Error())
		}
		return 0, err
	}
	return id, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*Provider, error) {
	prov, err := s.repository.Get(ctx, id)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	return prov, nil
}

func (s *Service) Update(ctx context.Context, prov *Provider) (uint64, error) {
	id, err := s.repository.Update(ctx, prov)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return 0, errors.ErrConflict.WithMsgf(err.Error())
		}
		if errors.As(err, new(NotFoundError)) {
			return 0, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return 0, err
	}
	return id, nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	return s.repository.Delete(ctx, id)
}
