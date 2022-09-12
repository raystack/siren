package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/odpf/siren/core/template"
)

type Template struct {
	ID        uint64         `db:"id"`
	Name      string         `db:"name"`
	Body      string         `db:"body"`
	Tags      pq.StringArray `db:"tags"`
	Variables string         `db:"variables"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (tmp *Template) FromDomain(t template.Template) error {
	tmp.ID = t.ID
	tmp.CreatedAt = t.CreatedAt
	tmp.UpdatedAt = t.UpdatedAt
	tmp.Name = t.Name
	tmp.Tags = t.Tags
	tmp.Body = t.Body
	jsonString, err := json.Marshal(t.Variables)
	if err != nil {
		return err
	}
	tmp.Variables = string(jsonString)
	return nil
}

func (tmp *Template) ToDomain() (*template.Template, error) {
	if tmp == nil {
		return nil, errors.New("template model is nil")
	}
	var variables []template.Variable
	jsonBlob := []byte(tmp.Variables)
	err := json.Unmarshal(jsonBlob, &variables)
	if err != nil {
		return nil, err
	}
	return &template.Template{
		ID:        tmp.ID,
		Name:      tmp.Name,
		Body:      tmp.Body,
		Tags:      tmp.Tags,
		CreatedAt: tmp.CreatedAt,
		UpdatedAt: tmp.UpdatedAt,
		Variables: variables,
	}, nil
}
