package silence

import (
	"context"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) Create(ctx context.Context, sil Silence) (string, error) {
	if err := sil.Validate(); err != nil {
		return "", err
	}
	return s.repository.Create(ctx, sil)
}

func (s *Service) List(ctx context.Context, filter Filter) ([]Silence, error) {
	return s.repository.List(ctx, filter)
}

func (s *Service) Get(ctx context.Context, id string) (Silence, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repository.SoftDelete(ctx, id)
}
