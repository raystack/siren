package receiver

import (
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db *gorm.DB
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r Repository) List() ([]*Receiver, error) {
	panic("implement me")
}

func (r Repository) Create(receiver *Receiver) (*Receiver, error) {
	panic("implement me")
}

func (r Repository) Get(u uint64) (*Receiver, error) {
	panic("implement me")
}

func (r Repository) Update(receiver *Receiver) (*Receiver, error) {
	panic("implement me")
}

func (r Repository) Delete(u uint64) error {
	panic("implement me")
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Receiver{})
	if err != nil {
		return err
	}
	return nil
}
