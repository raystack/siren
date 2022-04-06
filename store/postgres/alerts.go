package postgres

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/store/model"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewAlertRepository returns repository struct
func NewAlertRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&model.Alert{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Create(alert *domain.Alert) (*domain.Alert, error) {
	alertInModelType := &model.Alert{}
	alertInModelType.FromDomain(alert)
	result := r.db.Create(alertInModelType)
	if result.Error != nil {
		return nil, result.Error
	}
	res := alertInModelType.ToDomain()
	return &res, nil
}

func (r Repository) Get(resourceName string, providerId, startTime, endTime uint64) ([]domain.Alert, error) {
	var alerts []model.Alert
	selectQuery := fmt.Sprintf("select * from alerts where resource_name = '%s' AND provider_id = '%d' AND triggered_at BETWEEN to_timestamp('%d') AND to_timestamp('%d')",
		resourceName, providerId, startTime, endTime)
	result := r.db.Raw(selectQuery).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	resultInDomainType := make([]domain.Alert, 0, len(alerts))
	for i := 0; i < len(alerts); i++ {
		resultInDomainType = append(resultInDomainType, alerts[i].ToDomain())
	}
	return resultInDomainType, nil
}
