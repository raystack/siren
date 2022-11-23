package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

const receiverInsertQuery = `
INSERT INTO receivers (name, type, labels, configurations, created_at, updated_at)
    VALUES ($1, $2, $3, $4, now(), now())
RETURNING *
`

const receiverUpdateQuery = `
UPDATE receivers SET name=$2, labels=$3, configurations=$4, updated_at=now()
WHERE id = $1
RETURNING *
`

const receiverDeleteQuery = `
DELETE from receivers where id=$1
`

var receiverListQueryBuilder = sq.Select(
	"id",
	"name",
	"type",
	"labels",
	"configurations",
	"created_at",
	"updated_at",
).From("receivers")

// ReceiverRepository talks to the store to read or insert data
type ReceiverRepository struct {
	client    *Client
	tableName string
}

// NewReceiverRepository returns repository struct
func NewReceiverRepository(client *Client) *ReceiverRepository {
	return &ReceiverRepository{client, "receivers"}
}

func (r ReceiverRepository) List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error) {
	var queryBuilder = receiverListQueryBuilder
	if len(flt.ReceiverIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"id": flt.ReceiverIDs})
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

	receiversDomain := []receiver.Receiver{}
	for rows.Next() {
		var receiverModel model.Receiver
		if err := rows.StructScan(&receiverModel); err != nil {
			return nil, err
		}
		receiversDomain = append(receiversDomain, *receiverModel.ToDomain())
	}

	return receiversDomain, nil
}

func (r ReceiverRepository) Create(ctx context.Context, rcv *receiver.Receiver) error {
	if rcv == nil {
		return errors.New("receiver domain is nil")
	}

	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "INSERT", r.tableName, map[string]string{
	// 	"db.statement": receiverInsertQuery,
	// })
	// defer span.End()

	receiverModel := new(model.Receiver)
	receiverModel.FromDomain(*rcv)

	var createdReceiver model.Receiver
	if err := r.client.db.QueryRowxContext(ctx, receiverInsertQuery,
		receiverModel.Name,
		receiverModel.Type,
		receiverModel.Labels,
		receiverModel.Configurations,
	).StructScan(&createdReceiver); err != nil {
		err := checkPostgresError(err)
		if errors.Is(err, errDuplicateKey) {
			return provider.ErrDuplicate
		}
		return err
	}

	*rcv = *createdReceiver.ToDomain()

	return nil
}

func (r ReceiverRepository) Get(ctx context.Context, id uint64) (*receiver.Receiver, error) {
	query, args, err := receiverListQueryBuilder.Where("id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "SELECT", r.tableName, map[string]string{
	// 	"db.statement": query,
	// })
	// defer span.End()

	var receiverModel model.Receiver
	if err := r.client.db.GetContext(ctx, &receiverModel, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, receiver.NotFoundError{ID: id}
		}
		return nil, err
	}

	return receiverModel.ToDomain(), nil
}

func (r ReceiverRepository) Update(ctx context.Context, rcv *receiver.Receiver) error {
	if rcv == nil {
		return errors.New("receiver domain is nil")
	}

	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "UPDATE", r.tableName, map[string]string{
	// 	"db.statement": receiverUpdateQuery,
	// })
	// defer span.End()

	receiverModel := new(model.Receiver)
	receiverModel.FromDomain(*rcv)

	var updatedReceiver model.Receiver
	if err := r.client.db.QueryRowxContext(ctx, receiverUpdateQuery,
		receiverModel.ID,
		receiverModel.Name,
		receiverModel.Labels,
		receiverModel.Configurations,
	).StructScan(&updatedReceiver); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return receiver.NotFoundError{ID: receiverModel.ID}
		}
		return err
	}

	*rcv = *updatedReceiver.ToDomain()

	return nil
}

func (r ReceiverRepository) Delete(ctx context.Context, id uint64) error {
	// ctx, span := r.client.postgresTracer.StartSpan(ctx, "DELETE", r.tableName, map[string]string{
	// 	"db.statement": receiverDeleteQuery,
	// })
	// defer span.End()

	if _, err := r.client.db.ExecContext(ctx, receiverDeleteQuery, id); err != nil {
		return err
	}
	return nil
}
