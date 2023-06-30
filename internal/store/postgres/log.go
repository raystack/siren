package postgres

import (
	"context"
	"fmt"

	"github.com/raystack/siren/core/log"
	"github.com/raystack/siren/internal/store/model"
	"github.com/raystack/siren/pkg/errors"
	"github.com/raystack/siren/pkg/pgc"
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
	rows, err := r.client.QueryxContext(ctx, pgc.OpSelectAll, notificationLogTableName, fmt.Sprintf(`
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
	rows, err := r.client.QueryxContext(ctx, pgc.OpSelectAll, notificationLogTableName, fmt.Sprintf(`
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

	res, err := r.client.NamedExecContext(ctx, pgc.OpInsert, notificationLogTableName, notificationLogInsertNamedQuery, nssModel)
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

// func (r *SubscriptionRepository) Get(ctx context.Context, id uint64) (*subscription.Subscription, error) {
// 	query, args, err := subscriptionListQueryBuilder.Where("id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var subscriptionModel model.Subscription
// 	if err := r.client.QueryRowxContext(ctx, pgc.OpSelect, r.tableName, query, args...).StructScan(&subscriptionModel); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, subscription.NotFoundError{ID: id}
// 		}
// 		return nil, err
// 	}

// 	return subscriptionModel.ToDomain(), nil
// }

// func (r *SubscriptionRepository) Update(ctx context.Context, sub *subscription.Subscription) error {
// 	if sub == nil {
// 		return errors.New("subscription domain is nil")
// 	}

// 	subscriptionModel := new(model.Subscription)
// 	subscriptionModel.FromDomain(*sub)

// 	var newSubscriptionModel model.Subscription
// 	if err := r.client.QueryRowxContext(ctx, pgc.OpUpdate, r.tableName, subscriptionUpdateQuery,
// 		subscriptionModel.ID,
// 		subscriptionModel.NamespaceID,
// 		subscriptionModel.URN,
// 		subscriptionModel.Receiver,
// 		subscriptionModel.Match,
// 	).StructScan(&newSubscriptionModel); err != nil {
// 		err := pgc.CheckError(err)
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return subscription.NotFoundError{ID: subscriptionModel.ID}
// 		}
// 		if errors.Is(err, pgc.ErrDuplicateKey) {
// 			return subscription.ErrDuplicate
// 		}
// 		if errors.Is(err, pgc.ErrForeignKeyViolation) {
// 			return subscription.ErrRelation
// 		}
// 		return err
// 	}

// 	*sub = *newSubscriptionModel.ToDomain()

// 	return nil
// }

// func (r *SubscriptionRepository) Delete(ctx context.Context, id uint64) error {
// 	if _, err := r.client.ExecContext(ctx, pgc.OpDelete, r.tableName, subscriptionDeleteQuery, id); err != nil {
// 		return err
// 	}
// 	return nil
// }
