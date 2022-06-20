package alert

import (
	"time"
)

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service Service) Migrate() error {
	return service.repository.Migrate()
}

func (service Service) Create(alerts *Alerts) ([]Alert, error) {
	result := make([]Alert, 0, len(alerts.Alerts))

	for i := 0; i < len(alerts.Alerts); i++ {
		err := service.repository.Create(&alerts.Alerts[i])
		if err != nil {
			return nil, err
		}
		result = append(result, alerts.Alerts[i])
	}
	return result, nil
}

func (service Service) Get(resourceName string, providerId, startTime, endTime uint64) ([]Alert, error) {
	if endTime == 0 {
		endTime = uint64(time.Now().Unix())
	}

	return service.repository.Get(resourceName, providerId, startTime, endTime)
}
