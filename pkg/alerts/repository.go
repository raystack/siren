package alerts

import (
	"fmt"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Alert{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Create(alert *Alert) (*Alert, error) {
	result := r.db.Create(alert)
	if result.Error != nil {
		return nil, result.Error
	}
	return alert, nil
}

func (r Repository) Get(resourceName string, providerId, startTime, endTime uint64) ([]Alert, error) {
	var alerts []Alert
	selectQuery := fmt.Sprintf("select * from alerts where resource_name = '%s' AND provider_id = '%d' AND triggered_at BETWEEN to_timestamp('%d') AND to_timestamp('%d')",
		resourceName, providerId, startTime, endTime)
	result := r.db.Raw(selectQuery).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	return alerts, nil
}
