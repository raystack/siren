package alert

import (
	"context"
	"time"

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

func (s *Service) Create(ctx context.Context, alerts []*Alert) ([]Alert, error) {
	result := make([]Alert, 0, len(alerts))

	for _, alrt := range alerts {
		newAlert, err := s.repository.Create(ctx, alrt)
		if err != nil {
			if errors.Is(err, ErrRelation) {
				return nil, errors.ErrNotFound.WithMsgf(err.Error())
			}
			return nil, err
		}
		result = append(result, *newAlert)
	}
	return result, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Alert, error) {
	if flt.EndTime == 0 {
		flt.EndTime = uint64(time.Now().Unix())
	}

	return s.repository.List(ctx, flt)
}
