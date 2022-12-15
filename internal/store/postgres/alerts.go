package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/pgc"
)

const alertInsertQuery = `
INSERT INTO alerts (provider_id, namespace_id, resource_name, metric_name, metric_value, severity, rule, triggered_at, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
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
	client    *pgc.Client
	tableName string
}

// NewAlertRepository returns repository struct
func NewAlertRepository(client *pgc.Client) *AlertRepository {
	return &AlertRepository{client, "alerts"}
}

func (r AlertRepository) Create(ctx context.Context, alrt *alert.Alert) error {
	if alrt == nil {
		return errors.New("alert domain is nil")
	}

	var alertModel model.Alert
	alertModel.FromDomain(*alrt)

	var newAlertModel model.Alert
	if err := r.client.QueryRowxContext(ctx, pgc.OpInsert, r.tableName, alertInsertQuery,
		alertModel.ProviderID,
		alertModel.NamespaceID,
		alertModel.ResourceName,
		alertModel.MetricName,
		alertModel.MetricValue,
		alertModel.Severity,
		alertModel.Rule,
		alertModel.TriggeredAt,
	).StructScan(&newAlertModel); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, pgc.ErrForeignKeyViolation) {
			return alert.ErrRelation
		}
		return err
	}

	return nil
}

func (r AlertRepository) List(ctx context.Context, flt alert.Filter) ([]alert.Alert, error) {
	var queryBuilder = alertListQueryBuilder
	if flt.NamespaceID != 0 {
		queryBuilder = queryBuilder.Where("namespace_id = ?", flt.NamespaceID)
	}
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

	rows, err := r.client.QueryxContext(ctx, pgc.OpSelectAll, r.tableName, query, args...)
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
