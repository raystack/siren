package template

// Service handles business logic
type Service struct {
	repository Repository
}

// NewService returns repository struct
func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (service Service) Upsert(template *Template) error {
	return service.repository.Upsert(template)
}

func (service Service) Index(tag string) ([]Template, error) {
	return service.repository.Index(tag)
}

func (service Service) GetByName(name string) (*Template, error) {
	return service.repository.GetByName(name)
}

func (service Service) Delete(name string) error {
	return service.repository.Delete(name)
}

func (service Service) Render(name string, body map[string]string) (string, error) {
	return service.repository.Render(name, body)
}
