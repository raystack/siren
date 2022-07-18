package postgres

import (
	"context"
	"fmt"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/internal/store/model"
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
		rcv, err := r.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		receivers = append(receivers, *rcv)
	}

	return receivers, nil
}

func (r ReceiverRepository) Create(ctx context.Context, receiver *receiver.Receiver) error {
	m := new(model.Receiver)
	if err := m.FromDomain(receiver); err != nil {
		return err
	}

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
	rcv, err := rcvModel.ToDomain()
	if err != nil {
		return nil, err
	}
	return rcv, nil
}

func (r ReceiverRepository) Update(ctx context.Context, rcv *receiver.Receiver) error {
	var m model.Receiver
	if err := m.FromDomain(rcv); err != nil {
		return err
	}
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
