package codeexchange

import (
	"fmt"

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

func (r Repository) Upsert(accessToken *AccessToken) error {
	var existingAccessToken AccessToken
	result := r.db.Where(fmt.Sprintf("workspace = '%s'", accessToken.Workspace)).Find(&existingAccessToken)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		result = r.db.Create(accessToken)
	} else {
		result = r.db.Where("id = ?", existingAccessToken.ID).Updates(accessToken)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&AccessToken{})
	if err != nil {
		return err
	}
	return nil
}
