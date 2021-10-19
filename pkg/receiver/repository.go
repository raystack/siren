package receiver

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db            *gorm.DB
}


// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) List() ([]*Receiver, error) {
	var receivers []*Receiver
	selectQuery := fmt.Sprintf("select * from receivers")
	result := r.db.Raw(selectQuery).Find(&receivers)
	if result.Error != nil {
		return nil, result.Error
	}

	return receivers, nil
}

func (r Repository) Create(receiver *Receiver) (*Receiver, error) {
	var newReceiver Receiver
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

func (r Repository) Get(id uint64) (*Receiver, error) {
	var receiver Receiver
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&receiver)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &receiver, nil
}

func (r Repository) Update(receiver *Receiver) (*Receiver, error) {
	var newReceiver, existingReceiver Receiver
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

func (r Repository) Delete(id uint64) error {
	var receiver Receiver
	result := r.db.Where("id = ?", id).Delete(&receiver)
	return result.Error
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Receiver{})
	if err != nil {
		return err
	}
	return nil
}
