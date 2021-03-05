package domain

import (
	"time"
)

type Variable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Template struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	Name      string     `json:"name"`
	Body      string     `json:"body"`
	Tags      []string   `json:"tags"`
	Variables []Variable `json:"variables"`
}
