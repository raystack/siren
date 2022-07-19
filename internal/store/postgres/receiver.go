package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/pkg/errors"
)

// ReceiverRepository talks to the store to read or insert data
type ReceiverRepository struct {
	client *Client
}

// NewReceiverRepository returns repository struct
func NewReceiverRepository(client *Client) *ReceiverRepository {
	return &ReceiverRepository{client}
}

func (r ReceiverRepository) List(ctx context.Context, flt receiver.Filter) ([]receiver.Receiver, error) {
	var models []*model.Receiver
	result := r.client.db.WithContext(ctx)

	if len(flt.ReceiverIDs) > 0 {
		result = result.Where("id IN ?", flt.ReceiverIDs)
	}

	result = result.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	var receivers []receiver.Receiver
	for _, r := range models {
		receivers = append(receivers, *r.ToDomain())
	}

	return receivers, nil
}

func (r ReceiverRepository) Create(ctx context.Context, rcv *receiver.Receiver) error {
	if rcv == nil {
		return errors.New("receiver domain is nil")
	}

	m := new(model.Receiver)
	m.FromDomain(rcv)

	result := r.client.db.WithContext(ctx).Create(m)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r ReceiverRepository) Get(ctx context.Context, id uint64) (*receiver.Receiver, error) {
	rcvModel := new(model.Receiver)
	result := r.client.db.Where(fmt.Sprintf("id = %d", id)).Find(rcvModel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, receiver.NotFoundError{ID: id}
	}

	return rcvModel.ToDomain(), nil
}

func (r ReceiverRepository) Update(ctx context.Context, rcv *receiver.Receiver) error {
	if rcv == nil {
		return errors.New("receiver domain is nil")
	}

	var m model.Receiver
	m.FromDomain(rcv)

	result := r.client.db.WithContext(ctx).Where("id = ?", m.ID).Updates(m)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return receiver.NotFoundError{ID: rcv.ID}
	}

	return nil
}

func (r ReceiverRepository) Delete(ctx context.Context, id uint64) error {
	var receiver model.Receiver
	result := r.client.db.WithContext(ctx).Where("id = ?", id).Delete(&receiver)
	return result.Error
}
