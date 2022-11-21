package postgresq

import (
	"context"
	"errors"
	"fmt"
	"strings"

	newrelic "github.com/newrelic/go-agent"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/plugins/queues/postgresq/migrations"
	"go.opencensus.io/trace"
)

const (
	MESSAGE_QUEUE_SCHEMA_NAME = "notification"
	MESSAGE_QUEUE_TABLE_NAME  = MESSAGE_QUEUE_SCHEMA_NAME + ".message_queue"
)

type Strategy string

const (
	StrategyDefault Strategy = "default"
	StrategyDLQ     Strategy = "dlq"
)

type Queue struct {
	logger   log.Logger
	dbc      *db.Client
	strategy Strategy
}

var (
	successCallbackQuery = fmt.Sprintf(`
UPDATE %s
SET updated_at = $1, status = $2, try_count = $3
WHERE id = $4
`, MESSAGE_QUEUE_TABLE_NAME)

	errorCallbackQuery = fmt.Sprintf(`
UPDATE %s
SET updated_at = $1, status = $2, try_count = $3, last_error = $4, retryable = $5
WHERE id = $6
`, MESSAGE_QUEUE_TABLE_NAME)

	queueEnqueueNamedQuery = fmt.Sprintf(`
INSERT INTO %s
	(id, status, receiver_type, configs, details, last_error, max_tries, try_count, retryable,
	expired_at, created_at, updated_at)
    VALUES (:id,:status,:receiver_type,:configs,:details,:last_error,:max_tries,:try_count,:retryable,:expired_at,:created_at,:updated_at)
`, MESSAGE_QUEUE_TABLE_NAME)
)

func getQueueDequeueQuery(batchSize int, receiverTypesList string) string {
	return fmt.Sprintf(`
UPDATE %s
SET status = '%s', updated_at = now()
WHERE id IN (
    SELECT id
    FROM %s
    WHERE status = '%s' AND (expired_at < now() OR expired_at IS NULL) AND try_count < max_tries %s
    ORDER BY expired_at
    FOR UPDATE SKIP LOCKED
    LIMIT %d
)
RETURNING *
`, MESSAGE_QUEUE_TABLE_NAME, notification.MessageStatusPending, MESSAGE_QUEUE_TABLE_NAME, notification.MessageStatusEnqueued, receiverTypesList, batchSize)
}

func getDLQDequeueQuery(batchSize int, receiverTypesList string) string {
	return fmt.Sprintf(`
UPDATE %s
SET status = '%s', updated_at = now()
WHERE id IN (
    SELECT id
    FROM %s
    WHERE status = '%s' AND (expired_at < now() OR expired_at IS NULL) AND try_count < max_tries AND retryable IS TRUE %s
    ORDER BY expired_at
    FOR UPDATE SKIP LOCKED
    LIMIT %d
)
RETURNING *
`, MESSAGE_QUEUE_TABLE_NAME, notification.MessageStatusPending, MESSAGE_QUEUE_TABLE_NAME, notification.MessageStatusFailed, receiverTypesList, batchSize)
}

// New creates a new queue instance
func New(logger log.Logger, dbConfig db.Config, opts ...QueueOption) (*Queue, error) {
	q := &Queue{
		logger:   logger,
		strategy: StrategyDefault,
	}
	dbClient, err := db.New(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating postgres queue client: %w", err)
	}
	q.dbc = dbClient

	// create schema if not exist
	_, err = q.dbc.ExecContext(context.Background(), fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", MESSAGE_QUEUE_SCHEMA_NAME))
	if err != nil {
		return nil, err
	}

	dbConfig.URL = dbConfig.URL + fmt.Sprintf("&search_path=%s", MESSAGE_QUEUE_SCHEMA_NAME)

	if err := db.RunMigrations(dbConfig, migrations.FS, migrations.ResourcePath); err != nil {
		return nil, fmt.Errorf("error migrating postgres queue: %w", err)
	}

	for _, opt := range opts {
		opt(q)
	}

	return q, nil
}

// Dequeue pop the queue based on specific filters (receiver types or batch size) and process the messages with handlerFn
// message left in pending state that has expired or been updated long time ago means there was a failure when transforming row into a struct
func (q *Queue) Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []notification.Message) error) error {
	messages := []notification.Message{}

	receiverTypesQuery := getFilterReceiverTypes(receiverTypes)

	var dequeueQuery string
	if q.strategy == StrategyDLQ {
		dequeueQuery = getDLQDequeueQuery(batchSize, receiverTypesQuery)
	} else {
		dequeueQuery = getQueueDequeueQuery(batchSize, receiverTypesQuery)
	}
	rows, err := q.dbc.QueryxContext(ctx, dequeueQuery)
	if err != nil {
		return err
	}
	for rows.Next() {
		msg := NotificationMessage{}
		if err := rows.StructScan(&msg); err != nil {
			q.logger.Error("failed to transform message row into struct", "strategy", q.strategy, "error", err)
			continue
		}
		messages = append(messages, msg.ToDomain())
	}

	if len(messages) == 0 {
		return notification.ErrNoMessage
	} else {
		q.logger.Debug(fmt.Sprintf("dequeued %d messages with batch size %d", len(messages), batchSize), "strategy", q.strategy)
		if err := handlerFn(ctx, messages); err != nil {
			return fmt.Errorf("error processing dequeued message: %w", err)
		}
	}

	return nil
}

// Enqueue pushes messages to the queue
func (q *Queue) Enqueue(ctx context.Context, ms ...notification.Message) error {
	nr := newrelic.DatastoreSegment{
		Product:    nrProductName,
		Collection: fmt.Sprintf("%s.%s", ns.Name, ns.Set),
		Operation:  "Fetch",
		StartTime:  newrelic.FromContext(ctx).StartSegmentNow(),
	}
	defer nr.End()

	span := s.startSpan(ctx, "GET", ns)
	defer span.End()

	messages := []NotificationMessage{}
	for _, m := range ms {
		message := &NotificationMessage{}
		message.FromDomain(m)
		messages = append(messages, *message)
	}

	res, err := q.dbc.NamedExecContext(ctx, queueEnqueueNamedQuery, messages)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected when enqueueing messages")
	}
	return nil
}

// SuccessCallback is a callback that will be called once the message is succesfully handled by handlerFn
func (q *Queue) SuccessCallback(ctx context.Context, ms notification.Message) error {
	q.logger.Debug("marking a message as published", "strategy", q.strategy, "id", ms.ID)
	res, err := q.dbc.ExecContext(ctx, successCallbackQuery, ms.UpdatedAt, ms.Status, ms.TryCount, ms.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected when marking row as published")
	}
	q.logger.Debug("marked a message as published", "strategy", q.strategy, "id", ms.ID)
	return nil
}

// ErrorCallback is a callback that will be called once the message is failed to be handled by handlerFn
func (q *Queue) ErrorCallback(ctx context.Context, ms notification.Message) error {
	q.logger.Debug("marking a message as failed with", "strategy", q.strategy, "id", ms.ID)
	res, err := q.dbc.ExecContext(ctx, errorCallbackQuery, ms.UpdatedAt, ms.Status, ms.TryCount, ms.LastError, ms.Retryable, ms.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected when marking row as failed")
	}
	q.logger.Debug("marked a message as failed with", "strategy", q.strategy, "id", ms.ID)
	return nil
}

// Stop will close the db
func (q *Queue) Stop(ctx context.Context) error {
	return q.dbc.Close()
}

func getFilterReceiverTypes(receiverTypes []string) string {
	var receiverTypesQuery = ""
	if len(receiverTypes) > 0 {
		receiverTypesQuery = "AND receiver_type IN ("
		for _, rs := range receiverTypes {
			receiverTypesQuery += "'"
			receiverTypesQuery += rs
			receiverTypesQuery += "'"
			receiverTypesQuery += ","
		}
		receiverTypesQuery = strings.TrimSuffix(receiverTypesQuery, ",")
		receiverTypesQuery += ")"
	}
	return receiverTypesQuery
}

func (q *Queue) startSpan(ctx context.Context, query string) *trace.Span {
	// Refer https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/database.md
	_, span := trace.StartSpan(ctx, query)
	span.AddAttributes(
		trace.StringAttribute("db.system", "postgresql"),
		trace.StringAttribute("db.name", ns.Name),
		trace.StringAttribute("db.statement", ns.Set),
	)
	return span
}
