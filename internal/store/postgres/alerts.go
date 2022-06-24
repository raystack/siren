package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
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

func (r AlertRepository) Create(ctx context.Context, alrt *alert.Alert) (*alert.Alert, error) {
	var alertModel model.Alert
	if err := alertModel.FromDomain(alrt); err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).Create(&alertModel)
	if result.Error != nil {
		err := checkPostgresError(result.Error)
		if errors.Is(err, errForeignKeyViolation) {
			return nil, alert.ErrRelation
		}
		return nil, result.Error
	}

	newAlert, err := alertModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return newAlert, nil
}

func (r AlertRepository) List(ctx context.Context, flt alert.Filter) ([]alert.Alert, error) {
	var alertsModel []model.Alert
	selectQuery := fmt.Sprintf("select * from alerts where resource_name = '%s' AND provider_id = '%d' AND triggered_at BETWEEN to_timestamp('%d') AND to_timestamp('%d')",
		flt.ResourceName, flt.ProviderID, flt.StartTime, flt.EndTime)
	result := r.db.WithContext(ctx).Raw(selectQuery).Find(&alertsModel)
	if result.Error != nil {
		return nil, result.Error
	}

	var alerts []alert.Alert
	for _, am := range alertsModel {
		ad, err := am.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		alerts = append(alerts, *ad)
	}

	return alerts, nil
}
