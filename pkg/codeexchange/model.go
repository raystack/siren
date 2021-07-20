package codeexchange

import (
	"time"
)

type AccessToken struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessToken string
	Workspace   string
}

type ExchangeRepository interface {
	Upsert(*AccessToken) error
	Get(string) (string, error)
	Migrate() error
}
