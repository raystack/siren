package postgres

import (
	"errors"
	"fmt"
	"github.com/odpf/siren/store/model"
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

func (r ReceiverRepository) List() ([]*model.Receiver, error) {
	var receivers []*model.Receiver
	selectQuery := "select * from receivers"
	result := r.db.Raw(selectQuery).Find(&receivers)
	if result.Error != nil {
		return nil, result.Error
	}

	return receivers, nil
}

func (r ReceiverRepository) Create(receiver *model.Receiver) (*model.Receiver, error) {
	var newReceiver model.Receiver
	result := r.db.Create(receiver)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.db.Where(fmt.Sprintf("id = %d", receiver.Id)).Find(&newReceiver)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newReceiver, nil
}

func (r ReceiverRepository) Get(id uint64) (*model.Receiver, error) {
	var receiver model.Receiver
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&receiver)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("receiver not found: %d", id)
	}

	return &receiver, nil
}

func (r ReceiverRepository) Update(receiver *model.Receiver) (*model.Receiver, error) {
	var newReceiver, existingReceiver model.Receiver
	result := r.db.Where(fmt.Sprintf("id = %d", receiver.Id)).Find(&existingReceiver)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("receiver doesn't exist")
	} else {
		result = r.db.Where("id = ?", receiver.Id).Updates(receiver)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	result = r.db.Where(fmt.Sprintf("id = %d", receiver.Id)).Find(&newReceiver)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newReceiver, nil
}

func (r ReceiverRepository) Delete(id uint64) error {
	var receiver model.Receiver
	result := r.db.Where("id = ?", id).Delete(&receiver)
	return result.Error
}

func (r ReceiverRepository) Migrate() error {
	err := r.db.AutoMigrate(&model.Receiver{})
	if err != nil {
		return err
	}
	return nil
}
