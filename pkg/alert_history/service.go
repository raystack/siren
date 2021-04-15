package alert_history

import (
	"github.com/odpf/siren/domain"
	"gorm.io/gorm"
	"time"
)

// Service handles business logic
type Service struct {
	repository AlertHistoryRepository
}

// NewService returns repository struct
func NewService(db *gorm.DB) domain.AlertHistoryService {
	return &Service{NewRepository(db)}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Create(alerts *domain.Alerts) ([]domain.AlertHistoryObject, error) {
	result := make([]domain.AlertHistoryObject, 0, len(alerts.Alerts))
	for i := 0; i < len(alerts.Alerts); i++ {
		alertHistoryObject := &Alert{}
		err := alertHistoryObject.fromDomain(&alerts.Alerts[i])
		if err != nil {
			return nil, err
		}
		//err = isValid(alert)
		res, err := service.repository.Create(alertHistoryObject)
		if err != nil {
			return nil, err
		}
		createdAlertHistoryObj, err := res.toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, createdAlertHistoryObj)
	}
	return result, nil
}

func (service Service) Get(resource string, startTime uint32, endTime uint32) ([]domain.AlertHistoryObject, error) {
	if endTime == 0 {
		endTime = uint32(time.Now().Unix())
	}
	filteredAlerts, err := service.repository.Get(resource, startTime, endTime)
	if err != nil {
		return nil, err
	}
	result := make([]domain.AlertHistoryObject, 0, len(filteredAlerts))
	for i := 0; i < len(filteredAlerts); i++ {
		alertHistoryObj, err := filteredAlerts[i].toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, alertHistoryObj)
	}
	return result, nil
}
