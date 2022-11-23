package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const alertInsertQuery = `
INSERT INTO alerts (provider_id, resource_name, metric_name, metric_value, severity, rule, triggered_at, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *
`

var alertListQueryBuilder = sq.Select(
	"id",
	"provider_id",
	"resource_name",
	"metric_name",
	"metric_value",
	"severity",
	"rule",
	"triggered_at",
	"created_at",
	"updated_at",
).From("alerts")

// AlertRepository talks to the store to read or insert data
type AlertRepository struct {
	client    *Client
	tableName string
}

// NewAlertRepository returns repository struct
func NewAlertRepository(client *Client) *AlertRepository {
	return &AlertRepository{client, "alerts"}
}

func (r AlertRepository) Create(ctx context.Context, alrt *alert.Alert) error {
	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "INSERT", r.tableName, map[string]string{
	// 	"db.statement": alertInsertQuery,
	// })
	// defer span.End()

	if alrt == nil {
		return errors.New("alert domain is nil")
	}

	var alertModel model.Alert
	alertModel.FromDomain(*alrt)

	var newAlertModel model.Alert
	if err := r.client.db.QueryRowxContext(ctx, alertInsertQuery,
		alertModel.ProviderID,
		alertModel.ResourceName,
		alertModel.MetricName,
		alertModel.MetricValue,
		alertModel.Severity,
		alertModel.Rule,
		alertModel.TriggeredAt,
		alertModel.CreatedAt,
		alertModel.UpdatedAt,
	).StructScan(&newAlertModel); err != nil {
		err := checkPostgresError(err)
		if errors.Is(err, errForeignKeyViolation) {
			return alert.ErrRelation
		}
		return err
	}

	return nil
}

func (r AlertRepository) List(ctx context.Context, flt alert.Filter) ([]alert.Alert, error) {
	var queryBuilder = alertListQueryBuilder
	if flt.ResourceName != "" {
		queryBuilder = queryBuilder.Where("resource_name = ?", flt.ResourceName)
	}
	if flt.ProviderID != 0 {
		queryBuilder = queryBuilder.Where("provider_id = ?", flt.ProviderID)
	}

	if flt.StartTime != 0 && flt.EndTime != 0 {
		startTime := time.Unix(flt.StartTime, 0)
		endTime := time.Unix(flt.EndTime, 0)
		queryBuilder = queryBuilder.Where(sq.Expr("triggered_at BETWEEN ? AND ?", startTime, endTime))
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "SELECT_ALL", r.tableName, map[string]string{
	// 	"db.statement": query,
	// })
	// defer span.End()

	rows, err := r.client.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alertsDomain := []alert.Alert{}
	for rows.Next() {
		var alertModel model.Alert
		if err := rows.StructScan(&alertModel); err != nil {
			return nil, err
		}
		alertsDomain = append(alertsDomain, *alertModel.ToDomain())
	}

	return alertsDomain, nil
}
