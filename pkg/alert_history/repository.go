package alert_history

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

func (r Repository) Get(resource string, startTime uint32, endTime uint32) ([]Alert, error) {
	var alerts []Alert
	selectQuery := fmt.Sprintf("select * from alerts where resource = '%s' AND created_at BETWEEN to_timestamp('%d') AND to_timestamp('%d')",
		resource, startTime, endTime)
	result := r.db.Raw(selectQuery).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return alerts, nil
}
