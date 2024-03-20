package postgres

import (
	"context"
	"fmt"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/pgc"
	"go.nhat.io/otelsql"
	"go.opentelemetry.io/otel/attribute"
)

const notificationInsertQuery = `
INSERT INTO notifications (namespace_id, type, data, labels, valid_duration, template, unique_key, receiver_selectors, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
RETURNING *
`

// NotificationRepository talks to the store to read or insert data
type NotificationRepository struct {
	client *pgc.Client
}

// NewNotificationRepository returns NotificationRepository struct
func NewNotificationRepository(client *pgc.Client) *NotificationRepository {
	return &NotificationRepository{
		client: client,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, n notification.Notification) (notification.Notification, error) {
	nModel := new(model.Notification)
	nModel.FromDomain(n)

	var newNModel model.Notification

	// Instrumentation attributes
	attrs := []attribute.KeyValue{
		attribute.String("db.method", "Insert"),
		attribute.String("db.sql.table", "notifications"),
	}
	
	if err := r.client.QueryRowxContext(otelsql.AddMeterLabels(ctx, attrs...), notificationInsertQuery,
		nModel.NamespaceID,
		nModel.Type,
		nModel.Data,
		nModel.Labels,
		nModel.ValidDuration,
		nModel.Template,
		nModel.UniqueKey,
		nModel.ReceiverSelectors,
	).StructScan(&newNModel); err != nil {
		return notification.Notification{}, err
	}

	return newNModel.ToDomain(), nil
}

func (r *NotificationRepository) WithTransaction(ctx context.Context) context.Context {
	return r.client.WithTransaction(ctx, nil)
}

func (r *NotificationRepository) Rollback(ctx context.Context, err error) error {
	if txErr := r.client.Rollback(ctx); txErr != nil {
		return fmt.Errorf("rollback error %s with error: %w", txErr.Error(), err)
	}
	return nil
}

func (r *NotificationRepository) Commit(ctx context.Context) error {
	return r.client.Commit(ctx)
}
