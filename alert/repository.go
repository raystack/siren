package alert

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func (r Repository) UpsertSlack(credential SlackCredential) error {
	//r.db.Clauses(clause.OnConflict{DoUpdates: true})
	result := r.db.Create(&credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r Repository) UpsertPagerduty(credential PagerdutyCredential) error {
	result := r.db.Create(&credential)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}
