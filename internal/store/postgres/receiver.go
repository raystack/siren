package postgres

import (
	"fmt"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/internal/store/model"
	"gorm.io/gorm"
)

// ReceiverRepository talks to the store to read or insert data
type ReceiverRepository struct {
	db *gorm.DB
}

// NewReceiverRepository returns repository struct
func NewReceiverRepository(db *gorm.DB) *ReceiverRepository {
	return &ReceiverRepository{db}
}

func (r ReceiverRepository) List() ([]*receiver.Receiver, error) {
	var models []*model.Receiver
	selectQuery := "select * from receivers"
	result := r.db.Raw(selectQuery).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	var receivers []*receiver.Receiver
	for _, r := range models {
		rcv, err := r.ToDomain()
		if err != nil {
			// TODO log here
			continue
		}
		receivers = append(receivers, rcv)
	}

	return receivers, nil
}

func (r ReceiverRepository) Create(receiver *receiver.Receiver) error {
	m := new(model.Receiver)
	if err := m.FromDomain(receiver); err != nil {
		return err
	}

	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	newReceiver, err := m.ToDomain()
	if err != nil {
		return err
	}
	*receiver = *newReceiver
	return nil
}

func (r ReceiverRepository) Get(id uint64) (*receiver.Receiver, error) {
	rcvModel := new(model.Receiver)
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(rcvModel)
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

func (r ReceiverRepository) Update(rcv *receiver.Receiver) error {
	var m model.Receiver
	if err := m.FromDomain(rcv); err != nil {
		return err
	}
	result := r.db.Where("id = ?", m.ID).Updates(m)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return receiver.NotFoundError{ID: rcv.ID}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", m.ID)).Find(&m)
	if result.Error != nil {
		return result.Error
	}

	newRcv, err := m.ToDomain()
	if err != nil {
		return err
	}
	*rcv = *newRcv
	return nil
}

func (r ReceiverRepository) Delete(id uint64) error {
	var receiver model.Receiver
	result := r.db.Where("id = ?", id).Delete(&receiver)
	return result.Error
}
