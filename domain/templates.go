package domain

import (
	"time"
)

type Template struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Name      string    `json:"name"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
}

// TemplatesRepository interface
type TemplatesRepository interface {
	Upsert(*Template) (*Template, error)
}

// TemplatesService interface
type TemplatesService interface {
	Upsert(*Template) (*Template, error)
}
