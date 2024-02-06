package postgres

import (
	"context"
	"fmt"

	"github.com/goto/siren/core/log"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
)

const notificationLogTableName = "notification_log"

const notificationLogInsertNamedQuery = `
INSERT INTO notification_log
	(namespace_id, notification_id, subscription_id, alert_ids, receiver_id, silence_ids, created_at)
    VALUES (:namespace_id, :notification_id, :subscription_id, :alert_ids, :receiver_id, :silence_ids, now())
`

// LogRepository talks to the store to read or insert data
type LogRepository struct {
	client *pgc.Client
}

// NewLogRepository returns LogRepository struct
func NewLogRepository(client *pgc.Client) *LogRepository {
	return &LogRepository{
		client: client,
	}
}

func (r *LogRepository) ListAlertIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error) {
	rows, err := r.client.QueryxContext(ctx, fmt.Sprintf(`
	SELECT distinct unnest(alert_ids) AS alert_ids FROM %s WHERE silence_ids @> '{%s}'
	`, notificationLogTableName, silenceID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alertIDs []int64
	for rows.Next() {
		var alertID int64
		if err := rows.Scan(&alertID); err != nil {
			return nil, err
		}

		alertIDs = append(alertIDs, alertID)
	}

	return alertIDs, nil
}

func (r *LogRepository) ListSubscriptionIDsBySilenceID(ctx context.Context, silenceID string) ([]int64, error) {
	rows, err := r.client.QueryxContext(ctx, fmt.Sprintf(`
	SELECT distinct subscription_id FROM %s WHERE silence_ids @> '{%s}'
	`, notificationLogTableName, silenceID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptionIDs []int64
	for rows.Next() {
		var subscriptionID int64
		if err := rows.Scan(&subscriptionID); err != nil {
			return nil, err
		}

		subscriptionIDs = append(subscriptionIDs, subscriptionID)
	}

	return subscriptionIDs, nil
}

func (r *LogRepository) BulkCreate(ctx context.Context, nss []log.Notification) error {
	nssModel := []model.NotificationLog{}
	for _, ns := range nss {
		nsModel := new(model.NotificationLog)
		nsModel.FromDomain(ns)

		nssModel = append(nssModel, *nsModel)
	}

	res, err := r.client.NamedExecContext(ctx, notificationLogInsertNamedQuery, nssModel)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected when inserting notification subscribers")
	}
	return nil
}
