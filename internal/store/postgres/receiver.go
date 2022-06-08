package postgres

import (
	"errors"
	"fmt"

	"github.com/odpf/siren/domain"
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

func (r ReceiverRepository) List() ([]*domain.Receiver, error) {
	var models []*model.Receiver
	selectQuery := "select * from receivers"
	result := r.db.Raw(selectQuery).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	var receivers []*domain.Receiver
	for _, r := range models {
		receivers = append(receivers, r.ToDomain())
	}

	return receivers, nil
}

func (r ReceiverRepository) Create(receiver *domain.Receiver) error {
	m := new(model.Receiver)
	m.FromDomain(receiver)

	result := r.db.Create(m)
	if result.Error != nil {
		return result.Error
	}

	newReceiver := m.ToDomain()
	*receiver = *newReceiver
	return nil
}

func (r ReceiverRepository) Get(id uint64) (*domain.Receiver, error) {
	receiver := new(model.Receiver)
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(receiver)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("receiver not found: %d", id)
	}

	return receiver.ToDomain(), nil
}

func (r ReceiverRepository) Update(receiver *domain.Receiver) error {
	var m model.Receiver
	m.FromDomain(receiver)
	result := r.db.Where("id = ?", m.Id).Updates(m)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errors.New("receiver doesn't exist")
	}

	result = r.db.Where(fmt.Sprintf("id = %d", m.Id)).Find(&m)
	if result.Error != nil {
		return result.Error
	}

	newReceiver := m.ToDomain()
	*receiver = *newReceiver
	return nil
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
