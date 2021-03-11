package templates

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestService_GetByName(t *testing.T) {
	t.Run("should call repository GetByName method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := &domain.Template{
			ID: 1, Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		modelTemplate := &Template{
			ID: 1, Name: "foo", Body: "bar",
			Tags:      []string{"test"},
			Variables: `[{"name":"test-name", "default": "test-default", "description": "test-description", "type": "test-type"}]`,
		}
		repositoryMock.On("GetByName", "foo").Return(modelTemplate, nil).Once()
		result, err := dummyService.GetByName("foo")
		assert.Nil(t, err)
		assert.Equal(t, dummyTemplate, result)
		repositoryMock.AssertCalled(t, "GetByName", "foo")
	})

	t.Run("should call repository GetByName method and return nil if template not found", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("GetByName", "foo").Return(nil, nil).Once()
		result, err := dummyService.GetByName("foo")
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("should call repository GetByName method and return error if any", func(t *testing.T) {
		expectedError := errors.New("unexpected error")
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("GetByName", "foo").Return(nil, expectedError).Once()
		result, err := dummyService.GetByName("foo")
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("should call repository Delete method and return result", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", "foo").Return(nil).Once()
		result := dummyService.Delete("foo")
		assert.Nil(t, result)
		repositoryMock.AssertCalled(t, "Delete", "foo")
	})

	t.Run("should call repository Delete method and return error", func(t *testing.T) {
		expectedError := errors.New("unexpected error")
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", "foo").Return(expectedError).Once()
		result := dummyService.Delete("foo")
		assert.EqualError(t, result, expectedError.Error())
	})
}

func TestService_Index(t *testing.T) {
	t.Run("should call repository Index method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := []domain.Template{{
			ID: 1, Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}}
		modelTemplates := []Template{{
			ID: 1, Name: "foo", Body: "bar",
			Tags:      []string{"test"},
			Variables: `[{"name":"test-name", "default": "test-default", "description": "test-description", "type": "test-type"}]`,
		}}
		repositoryMock.On("Index", "foo").Return(modelTemplates, nil).Once()
		result, err := dummyService.Index("foo")
		assert.Nil(t, err)
		assert.Equal(t, dummyTemplate, result)
		repositoryMock.AssertCalled(t, "Index", "foo")
	})

	t.Run("should call repository Index method and return error if any", func(t *testing.T) {
		expectedError := errors.New("unexpected error")
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Index", "foo").Return(nil, expectedError).Once()
		result, err := dummyService.Index("foo")
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
	})
}

func TestService_Upsert(t *testing.T) {
	t.Run("should perform name validation", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := &domain.Template{
			Name: "", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		result, err := dummyService.Upsert(dummyTemplate)
		assert.EqualError(t, err, "name cannot be empty")
		assert.Nil(t, result)
		repositoryMock.AssertNotCalled(t, "UpsertSlack")
	})

	t.Run("should perform body validation", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		result, err := dummyService.Upsert(dummyTemplate)
		assert.EqualError(t, err, "body cannot be empty")
		assert.Nil(t, result)
		repositoryMock.AssertNotCalled(t, "UpsertSlack")
	})

	t.Run("should call repository UpsertSlack method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		modelTemplate := &Template{
			ID: 0, Name: "foo", Body: "bar",
			Tags:      []string{"test"},
			Variables: `[{"name":"test-name","type":"test-type","default":"test-default","description":"test-description"}]`,
		}
		repositoryMock.On("Upsert", modelTemplate).Return(modelTemplate, nil).Once()
		result, err := dummyService.Upsert(dummyTemplate)
		assert.Nil(t, err)
		assert.Equal(t, dummyTemplate, result)
		repositoryMock.AssertCalled(t, "Upsert", modelTemplate)
	})

	t.Run("should call repository Upsert method and return error if any", func(t *testing.T) {
		expectedError := errors.New("unexpected error")
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		modelTemplate := &Template{
			ID: 0, Name: "foo", Body: "bar",
			Tags:      []string{"test"},
			Variables: `[{"name":"test-name","type":"test-type","default":"test-default","description":"test-description"}]`,
		}
		repositoryMock.On("Upsert", modelTemplate).Return(nil, expectedError).Once()
		result, err := dummyService.Upsert(dummyTemplate)
		assert.Nil(t, result)
		assert.EqualError(t, err, expectedError.Error())
	})
}

func TestService_Render(t *testing.T) {
	t.Run("should call repository Render method and return result", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		input := make(map[string]string)
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Render", "foo", input).Return("foo bar baz", nil).Once()
		result, err := dummyService.Render("foo", input)
		assert.Nil(t, err)
		assert.Equal(t, "foo bar baz", result)
		repositoryMock.AssertCalled(t, "Render", "foo", input)
	})

	t.Run("should call repository Render method and return error", func(t *testing.T) {
		expectedError := errors.New("unexpected error")
		repositoryMock := &TemplatesRepositoryMock{}
		input := make(map[string]string)
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Render", "foo", input).Return("", expectedError).Once()
		result, err := dummyService.Render("foo", input)
		assert.Empty(t, result)
		assert.EqualError(t, err, expectedError.Error())
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &TemplatesRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
