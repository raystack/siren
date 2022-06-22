package template

import "github.com/odpf/siren/pkg/errors"

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service Service) Upsert(template *Template) error {
	if err := service.repository.Upsert(template); err != nil {
		if errors.Is(err, ErrDuplicate) {
			return errors.ErrConflict.WithMsgf(err.Error())
		}
		return err
	}
	return nil
}

func (service Service) Index(tag string) ([]Template, error) {
	return service.repository.Index(tag)
}

func (service Service) GetByName(name string) (*Template, error) {
	tmpl, err := service.repository.GetByName(name)
	if err != nil {
		if errors.As(err, new(NotFoundError)) {
			return nil, errors.ErrNotFound.WithMsgf(err.Error())
		}
		return nil, err
	}
	return tmpl, nil
}

func (service Service) Delete(name string) error {
	return service.repository.Delete(name)
}

func (service Service) Render(name string, body map[string]string) (string, error) {
	return service.repository.Render(name, body)
}
