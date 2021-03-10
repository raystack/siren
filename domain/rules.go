package domain

import "time"

type RuleVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Rule struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	Name      string         `json:"name"`
	Namespace string         `json:"namespace"`
	Entity    string         `json:"entity"`
	GroupName string         `json:"group_name"`
	Template  string         `json:"template"`
	Status    string         `json:"status"`
	Variables []RuleVariable `json:"variables"`
}

// RuleService interface
type RuleService interface {
	Upsert(*Rule) (*Rule, error)
	Get(string) ([]Rule, error)
	Migrate() error
}
