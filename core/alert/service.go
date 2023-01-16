package alert

import (
	"context"
	"time"

	"github.com/odpf/siren/pkg/errors"
)

// Service handles business logic
type Service struct {
	repository Repository
	registry   map[string]AlertTransformer
}

// NewService returns repository struct
func NewService(repository Repository, registry map[string]AlertTransformer) *Service {
	return &Service{repository, registry}
}

func (s *Service) CreateAlerts(ctx context.Context, providerType string, providerID uint64, namespaceID uint64, body map[string]interface{}) ([]Alert, int, error) {
	pluginService, err := s.getProviderPluginService(providerType)
	if err != nil {
		return nil, 0, err
	}

	alerts, firingLen, err := pluginService.TransformToAlerts(ctx, providerID, namespaceID, body)
	if err != nil {
		return nil, 0, err
	}

	for _, alrt := range alerts {
		createdAlert, err := s.repository.Create(ctx, alrt)
		if err != nil {
			if errors.Is(err, ErrRelation) {
				return nil, 0, errors.ErrNotFound.WithMsgf(err.Error())
			}
			return nil, 0, err
		}
		alrt.ID = createdAlert.ID
	}

	return alerts, firingLen, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Alert, error) {
	if flt.EndTime == 0 {
		flt.EndTime = time.Now().Unix()
	}

	return s.repository.List(ctx, flt)
}

func (s *Service) getProviderPluginService(providerType string) (AlertTransformer, error) {
	pluginService, exist := s.registry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}
