package domain

import "time"

type RuleVariable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}

type Rule struct {
	Id                uint64         `json:"id"`
	Name              string         `json:"name"`
	Enabled           bool           `json:"enabled" validate:"required"`
	GroupName         string         `json:"group_name" validate:"required"`
	Namespace         string         `json:"namespace" validate:"required"`
	Template          string         `json:"template" validate:"required"`
	Variables         []RuleVariable `json:"variables" validate:"required,dive,required"`
	ProviderNamespace uint64         `json:"provider_namespace" validate:"required"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// RuleService interface
type RuleService interface {
	Upsert(*Rule) (*Rule, error)
	Get(string, string, string, string) ([]Rule, error)
	Migrate() error
}
