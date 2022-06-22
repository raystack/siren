package postgres

import (
	"fmt"

	"github.com/odpf/siren/core/alert"
	"gorm.io/gorm"
)

// AlertRepository talks to the store to read or insert data
type AlertRepository struct {
	db *gorm.DB
}

// NewAlertRepository returns repository struct
func NewAlertRepository(db *gorm.DB) *AlertRepository {
	return &AlertRepository{db}
}

func (r AlertRepository) Create(alert *alert.Alert) error {
	result := r.db.Create(alert)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r AlertRepository) Get(resourceName string, providerId, startTime, endTime uint64) ([]alert.Alert, error) {
	var alerts []alert.Alert
	selectQuery := fmt.Sprintf("select * from alerts where resource_name = '%s' AND provider_id = '%d' AND triggered_at BETWEEN to_timestamp('%d') AND to_timestamp('%d')",
		resourceName, providerId, startTime, endTime)
	result := r.db.Raw(selectQuery).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	return alerts, nil
}
