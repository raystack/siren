package template

import (
	"bytes"
	"context"

	texttemplate "text/template"

	"github.com/odpf/siren/pkg/errors"
)

const (
	leftDelim  = "[["
	rightDelim = "]]"
)

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (s *Service) Upsert(ctx context.Context, template *Template) error {
	err := s.repository.Upsert(ctx, template)
	if err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		return err
	}
	return nil
}

func (s *Service) List(ctx context.Context, flt Filter) ([]Template, error) {
	return s.repository.List(ctx, flt)
}

func (s *Service) GetByName(ctx context.Context, name string) (*Template, error) {
	tmpl, err := s.repository.GetByName(ctx, name)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	return tmpl, nil
}

func (s *Service) Delete(ctx context.Context, name string) error {
	return s.repository.Delete(ctx, name)
}

func (s *Service) Render(ctx context.Context, name string, requestVariables map[string]string) (string, error) {
	templateFromDB, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return "", err
	}

	enrichedVariables := enrichWithDefaults(templateFromDB.Variables, requestVariables)
	var tpl bytes.Buffer
	tmpl, err := texttemplate.New("parser").Delims(leftDelim, rightDelim).Parse(templateFromDB.Body)
	if err != nil {
		return "", errors.ErrInvalid.WithMsgf("failed to parse template body").WithCausef(err.Error())
	}
	err = tmpl.Execute(&tpl, enrichedVariables)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func enrichWithDefaults(variables []Variable, requestVariables map[string]string) map[string]string {
	result := make(map[string]string)
	for _, variable := range variables {
		name := variable.Name
		defaultValue := variable.Default
		val, ok := requestVariables[name]
		if ok {
			result[name] = val
		} else {
			result[name] = defaultValue
		}
	}
	return result
}
