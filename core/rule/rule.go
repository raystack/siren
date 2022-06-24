package rule

import (
	"context"
	"time"
)

//go:generate mockery --name=Repository -r --case underscore --with-expecter --structname RuleRepository --filename rule_repository.go --output=./mocks
type Repository interface {
	UpsertWithTx(context.Context, *Rule, func() error) (uint64, error)
	List(context.Context, Filter) ([]Rule, error)
}

type RuleVariable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}

type Rule struct {
	ID                uint64         `json:"id"`
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
