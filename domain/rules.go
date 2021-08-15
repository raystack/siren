package domain

import (
	"time"
	"gopkg.in/go-playground/validator.v9"
)

type RuleVariable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type"`
	Value       string `json:"value" validate:"required"`
	Description string `json:"description"`
}

type Rule struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Name      string         `json:"name"`
	Namespace string         `json:"namespace" validate:"required"`
	Entity    string         `json:"entity" validate:"required"`
	GroupName string         `json:"group_name" validate:"required"`
	Template  string         `json:"template" validate:"required"`
	Status    string         `json:"status" validate:"required,statusChecker"`
	Variables []RuleVariable `json:"variables" validate:"required,dive,required"`
}

// RuleService interface
type RuleService interface {
	Upsert(*Rule) (*Rule, error)
	Get(string, string, string, string, string) ([]Rule, error)
	Migrate() error
}

func (rs *Rule) Validate() error {
	v := validator.New()
	_ = v.RegisterValidation("statusChecker", func(fl validator.FieldLevel) bool {
		return fl.Field().Interface().(string) == "enabled" || fl.Field().Interface().(string) == "disabled"
	})
	return v.Struct(rs)
}
