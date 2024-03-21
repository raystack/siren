package template

import (
	"bytes"
	"context"

	texttemplate "text/template"

	"github.com/goto/siren/pkg/errors"
)

const (
	defaultLeftDelim  = "[["
	defaultrightDelim = "]]"

	DelimMessageLeft  = "{{"
	DelimMessageRight = "}}"
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

// TODO might want to delete this and use the static function instead
func (s *Service) Render(ctx context.Context, name string, requestVariables map[string]string) (string, error) {
	templateFromDB, err := s.repository.GetByName(ctx, name)
	if err != nil {
		return "", err
	}

	return RenderWithEnrichedDefault(templateFromDB.Body, templateFromDB.Variables, requestVariables)
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

func RenderWithEnrichedDefault(templateBody string, templateVars []Variable, requestVariables map[string]string) (string, error) {
	enrichedVariables := enrichWithDefaults(templateVars, requestVariables)
	return RenderBody(templateBody, enrichedVariables, defaultLeftDelim, defaultrightDelim)
}

func RenderBody(templateBody string, aStruct interface{}, leftDelim, rightDelim string) (string, error) {
	var tpl bytes.Buffer
	tmpl, err := texttemplate.New("parser").Funcs(defaultFuncMap).Delims(leftDelim, rightDelim).Parse(templateBody)
	if err != nil {
		return "", errors.ErrInvalid.WithMsgf("failed to parse template body").WithMsgf(err.Error())
	}
	err = tmpl.Execute(&tpl, &aStruct)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
