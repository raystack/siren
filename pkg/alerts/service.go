package alerts

import (
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
	"time"
)

// Service handles business logic
type Service struct {
	repository AlertRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.AlertService {
	return &Service{NewRepository(db)}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Create(alerts *domain.Alerts) ([]domain.Alert, error) {
	result := make([]domain.Alert, 0, len(alerts.Alerts))

	for i := 0; i < len(alerts.Alerts); i++ {
		alertHistoryObject := &Alert{}
		alertHistoryObject.fromDomain(&alerts.Alerts[i])
		res, err := service.repository.Create(alertHistoryObject)
		if err != nil {
			return nil, err
		}
		createdAlertHistoryObj := res.toDomain()
		result = append(result, createdAlertHistoryObj)
	}
	return result, nil
}

func (service Service) Get(resourceName string, providerId, startTime, endTime uint64) ([]domain.Alert, error) {
	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}

	filteredAlerts, err := service.repository.Get(resourceName, providerId, startTime, endTime)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Alert, 0, len(filteredAlerts))
	for i := 0; i < len(filteredAlerts); i++ {
		alertHistoryObj := filteredAlerts[i].toDomain()
		result = append(result, alertHistoryObj)
	}
	return result, nil
}
