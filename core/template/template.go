package template

import (
	"context"
	"time"
)

const (
	ReservedName_SystemDefault = "system-default"
)

type Repository interface {
	Upsert(context.Context, *Template) error
	List(context.Context, Filter) ([]Template, error)
	GetByName(context.Context, string) (*Template, error)
	Delete(context.Context, string) error
}

type Variable struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type Template struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name" validate:"required"`
	Body      string     `json:"body" validate:"required"`
	Tags      []string   `json:"tags" validate:"required"`
	Variables []Variable `json:"variables" validate:"required,dive,required"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func IsReservedName(templateName string) bool {
	return (templateName == ReservedName_SystemDefault)
}

type FileBody struct {
	Rule    `yaml:",inline"`
	Message `yaml:",inline"`
}

type Message struct {
	ReceiverType string `yaml:"receiver_type,omitempty"` // mandatory
	Content      string `yaml:"content,omitempty"`
}

type Rule struct {
	Record      string            `yaml:"record,omitempty"`
	Alert       string            `yaml:"alert,omitempty"`
	Expr        string            `yaml:"expr,omitempty"` // mandatory
	For         string            `yaml:"for,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type File struct {
	Name       string     `yaml:"name"`
	ApiVersion string     `yaml:"apiVersion"`
	Type       string     `yaml:"type"`
	Body       []FileBody `yaml:"body"`
	Message    Message    `yaml:"message"`
	Tags       []string   `yaml:"tags"`
	Delete     bool       `yaml:"delete"`
	Variables  []Variable `yaml:"variables"`
}
