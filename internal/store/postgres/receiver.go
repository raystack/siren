package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/internal/store/model"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
)

const receiverInsertQuery = `
INSERT INTO receivers (name, type, labels, configurations, parent_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, now(), now())
RETURNING *
`

const receiverUpdateQuery = `
UPDATE receivers SET name=$2, labels=$3, configurations=$4, updated_at=now()
WHERE id = $1
RETURNING *
`

const receiverPatchLabelsQuery = `
UPDATE receivers SET labels=$2, updated_at=now()
WHERE id = $1
RETURNING *
`

const receiverDeleteQuery = `
DELETE from receivers where id=$1
`

var receiverListQueryBuilder = sq.Select(
	"r.id as id",
	"r.name as name",
	"r.type as type",
	"r.labels as labels",
	"r.parent_id as parent_id",
	"r.created_at as created_at",
	"r.updated_at as updated_at",
	"r.configurations as configurations",
).From("receivers r")

var receiverListLeftJoinSelfQueryBuilder = sq.Select(
	"r.id",
	"r.name",
	"r.type",
	"r.labels",
	"r.parent_id",
	"r.created_at",
	"r.updated_at",
).Column(
	sq.Expr("r.configurations || COALESCE(p.configurations, '{}'::jsonb) AS configurations"),
).From("receivers r")

// ReceiverRepository talks to the store to read or insert data
type ReceiverRepository struct {
	client *pgc.Client
}

// NewReceiverRepository returns repository struct
func NewReceiverRepository(client *pgc.Client) *ReceiverRepository {
	return &ReceiverRepository{client}
}

func (r ReceiverRepository) List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error) {
	var queryBuilder sq.SelectBuilder
	if flt.Expanded {
		queryBuilder = receiverListLeftJoinSelfQueryBuilder.LeftJoin("receivers p ON r.parent_id = p.id")
	} else {
		queryBuilder = receiverListQueryBuilder
	}

	if len(flt.ReceiverIDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"r.id": flt.ReceiverIDs})
	}

	// given map of string from input [lf], look for rows that [lf] exist in labels column in DB
	if len(flt.MultipleLabels) != 0 {
		var matchLabelsExpression = sq.Or{}
		for _, labels := range flt.MultipleLabels {
			labelsJSON, err := json.Marshal(labels)
			if err != nil {
				return nil, errors.ErrInvalid.WithCausef("problem marshalling labels %v json to string with err: %s", labels, err.Error())
			}
			matchLabelsExpression = append(
				matchLabelsExpression,
				sq.Expr(fmt.Sprintf("r.labels @> '%s'::jsonb", string(json.RawMessage(labelsJSON)))),
			)
		}
		queryBuilder = queryBuilder.Where(matchLabelsExpression)
	}

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.client.QueryxContext(ctx, query, args...)
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

	receiverModel := new(model.Receiver)
	receiverModel.FromDomain(*rcv)

	var createdReceiver model.Receiver
	if err := r.client.QueryRowxContext(ctx, receiverInsertQuery,
		receiverModel.Name,
		receiverModel.Type,
		receiverModel.Labels,
		receiverModel.Configurations,
		receiverModel.ParentID,
	).StructScan(&createdReceiver); err != nil {
		err = pgc.CheckError(err)
		if errors.Is(err, pgc.ErrDuplicateKey) {
			return provider.ErrDuplicate
		}
		return err
	}

	*rcv = *createdReceiver.ToDomain()

	return nil
}

func (r ReceiverRepository) Get(ctx context.Context, id uint64, flt receiver.Filter) (*receiver.Receiver, error) {

	var queryBuilder sq.SelectBuilder
	if flt.Expanded {
		queryBuilder = receiverListLeftJoinSelfQueryBuilder.LeftJoin("receivers p ON r.parent_id = p.id")
	} else {
		queryBuilder = receiverListQueryBuilder
	}

	query, args, err := queryBuilder.Where("r.id = ?", id).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var receiverModel model.Receiver
	if err := r.client.GetContext(ctx, &receiverModel, query, args...); err != nil {
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

	receiverModel := new(model.Receiver)
	receiverModel.FromDomain(*rcv)

	var updatedReceiver model.Receiver
	if err := r.client.QueryRowxContext(ctx, receiverUpdateQuery,
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

func (r ReceiverRepository) PatchLabels(ctx context.Context, rcv *receiver.Receiver) error {
	if rcv == nil {
		return errors.ErrInvalid.WithCausef("receiver cannot be nil")
	}

	receiverModel := new(model.Receiver)
	receiverModel.FromDomain(*rcv)

	var patchedLabelReceiver model.Receiver
	if err := r.client.QueryRowxContext(ctx, receiverPatchLabelsQuery,
		receiverModel.ID,
		receiverModel.Labels,
	).StructScan(&patchedLabelReceiver); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return receiver.NotFoundError{ID: receiverModel.ID}
		}
		return err
	}

	*rcv = *patchedLabelReceiver.ToDomain()

	return nil
}

func (r ReceiverRepository) Delete(ctx context.Context, id uint64) error {
	if _, err := r.client.ExecContext(ctx, receiverDeleteQuery, id); err != nil {
		return err
	}
	return nil
}
