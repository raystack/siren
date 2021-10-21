package rules

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var truebool = true

func TestService_Upsert(t *testing.T) {
	t.Run("should call repository Upsert method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockCortexClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := Service{repository: repositoryMock, client: mockCortexClient, templateService: mockTemplateService}
		dummyRule := &domain.Rule{
			Id: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}
		modelRule := &Rule{
			Id: 1, Name: "bar", Enabled: &truebool, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables:         `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
			ProviderNamespace: 1,
		}
		repositoryMock.On("Upsert", modelRule, mockCortexClient, mockTemplateService).Return(modelRule, nil).Once()
		result, err := dummyService.Upsert(dummyRule)
		assert.Nil(t, err)
		assert.Equal(t, dummyRule, result)
		repositoryMock.AssertCalled(t, "Upsert", modelRule, mockCortexClient, mockTemplateService)
	})

	t.Run("should call repository Upsert method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockCortexClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := Service{repository: repositoryMock, client: mockCortexClient, templateService: mockTemplateService}
		dummyRule := &domain.Rule{
			Id: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}
		modelRule := &Rule{
			Id: 1, Name: "bar", Enabled: &truebool, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables:         `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
			ProviderNamespace: 1,
		}
		repositoryMock.On("Upsert", modelRule, mockCortexClient, mockTemplateService).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.Upsert(dummyRule)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Upsert", modelRule, mockCortexClient, mockTemplateService)
	})

	t.Run("should call repository Upsert method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockCortexClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := Service{repository: repositoryMock, client: mockCortexClient, templateService: mockTemplateService}
		dummyRule := &domain.Rule{
			Id: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}
		modelRule := &Rule{
			Id: 1, Name: "bar", Enabled: &truebool, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables:         `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
			ProviderNamespace: 1,
		}
		repositoryMock.On("Upsert", modelRule, mockCortexClient, mockTemplateService).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.Upsert(dummyRule)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Upsert", modelRule, mockCortexClient, mockTemplateService)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyRules := []domain.Rule{{
			Id: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}}
		modelRules := []Rule{{
			Id: 1, Name: "bar", Enabled: &truebool, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables:         `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
			ProviderNamespace: 1,
		}}
		repositoryMock.On("Get", "foo", "gojek", "test-group", "test-tmpl").
			Return(modelRules, nil).Once()

		result, err := dummyService.Get("foo", "gojek", "test-group", "test-tmpl")
		assert.Nil(t, err)
		assert.Equal(t, dummyRules, result)
		repositoryMock.AssertCalled(t, "Get", "foo", "gojek", "test-group", "test-tmpl")
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", "foo", "", "", "").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.Get("foo", "", "", "")
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Get", "foo", "", "", "")
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
