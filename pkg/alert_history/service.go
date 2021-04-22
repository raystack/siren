package alert_history

import (
	"errors"
	"fmt"
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

// Creating all valid alert history objects from array of objects
func (service Service) Create(alerts *domain.Alerts) ([]domain.AlertHistoryObject, error) {
	result := make([]domain.AlertHistoryObject, 0, len(alerts.Alerts))
	badAlertHistoryObjectCount := 0
	for i := 0; i < len(alerts.Alerts); i++ {
		alertHistoryObject := &Alert{}
		alertHistoryObject.fromDomain(&alerts.Alerts[i])
		if !isValid(alertHistoryObject) {
			badAlertHistoryObjectCount++
			continue
		}
		res, err := service.repository.Create(alertHistoryObject)
		if err != nil {
			return nil, err
		}
		createdAlertHistoryObj := res.toDomain()
		result = append(result, createdAlertHistoryObj)
	}

	if badAlertHistoryObjectCount > 0 {
		return result,
			errors.New(fmt.Sprintf("alert history parameters missing for %d alerts", badAlertHistoryObjectCount))
	}
	return result, nil
}

func isValid(alert *Alert) bool {
	return !(alert.Resource == "" || alert.Template == "" ||
		alert.MetricValue == "" || alert.MetricName == "" ||
		alert.Level == "")
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
		alertHistoryObj := filteredAlerts[i].toDomain()
		result = append(result, alertHistoryObj)
	}
	return result, nil
}
