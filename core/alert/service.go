package alert

import (
	"context"
	"time"

	"github.com/raystack/siren/pkg/errors"
)

//go:generate mockery --name=LogService -r --case underscore --with-expecter --structname LogService --filename log_service.go --output=./mocks
type LogService interface {
	ListAlertIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error)
}

// Service handles business logic
type Service struct {
	repository Repository
	logService LogService
	registry   map[string]AlertTransformer
}

// NewService returns repository struct
func NewService(repository Repository, logService LogService, registry map[string]AlertTransformer) *Service {
	return &Service{repository, logService, registry}
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

	for i := 0; i < len(alerts); i++ {
		createdAlert, err := s.repository.Create(ctx, alerts[i])
		if err != nil {
			if errors.Is(err, ErrRelation) {
				return nil, 0, errors.ErrNotFound.WithMsgf(err.Error())
			}
			return nil, 0, err
		}
		alerts[i].ID = createdAlert.ID
	}

	return alerts, firingLen, nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Alert, error) {
	if flt.EndTime == 0 {
		flt.EndTime = time.Now().Unix()
	}

	if flt.SilenceID != "" {
		alertIDs, err := s.logService.ListAlertIDsBySilenceID(ctx, flt.SilenceID)
		if err != nil {
			return nil, err
		}
		flt.IDs = alertIDs
	}

	return s.repository.List(ctx, flt)
}

func (s *Service) UpdateSilenceStatus(ctx context.Context, alertIDs []int64, hasSilenced bool, hasNonSilenced bool) error {
	return s.repository.BulkUpdateSilence(ctx, alertIDs, silenceStatus(hasSilenced, hasNonSilenced))
}

func (s *Service) getProviderPluginService(providerType string) (AlertTransformer, error) {
	pluginService, exist := s.registry[providerType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported provider type: %q", providerType)
	}
	return pluginService, nil
}
