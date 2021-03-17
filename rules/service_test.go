package rules

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Upsert(t *testing.T) {
	t.Run("should call repository Upsert method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockCortexClient := &cortexCallerMock{}
		dummyService := Service{repository: repositoryMock, client: mockCortexClient}
		dummyRule := &domain.Rule{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		modelRule := &Rule{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
		}
		repositoryMock.On("Upsert", modelRule, mockCortexClient).Return(modelRule, nil).Once()
		result, err := dummyService.Upsert(dummyRule)
		assert.Nil(t, err)
		assert.Equal(t, dummyRule, result)
		repositoryMock.AssertCalled(t, "Upsert", modelRule, mockCortexClient)
	})

	t.Run("should call repository Upsert method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockCortexClient := &cortexCallerMock{}
		dummyService := Service{repository: repositoryMock, client: mockCortexClient}
		dummyRule := &domain.Rule{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}
		modelRule := &Rule{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: `[{"name":"test-name","type":"test-type","value":"test-value","description":"test-description"}]`,
		}
		repositoryMock.On("Upsert", modelRule, mockCortexClient).Return(nil, errors.New("random error")).Once()
		result, err := dummyService.Upsert(dummyRule)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Upsert", modelRule, mockCortexClient)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		dummyRules := []domain.Rule{{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			},
		}}
		modelRules := []Rule{{
			ID: 1, Name: "bar", Namespace: "baz",
			Entity: "gojek", GroupName: "test-group", Template: "test-tmpl", Status: "enabled",
			Variables: `[{"name":"test-name", "value": "test-value", "description": "test-description", "type": "test-type"}]`,
		}}
		repositoryMock.On("Get", "foo", "gojek", "test-group", "enabled", "test-tmpl").Return(modelRules, nil).Once()
		result, err := dummyService.Get("foo", "gojek", "test-group", "enabled", "test-tmpl")
		assert.Nil(t, err)
		assert.Equal(t, dummyRules, result)
		repositoryMock.AssertCalled(t, "Get", "foo", "gojek", "test-group", "enabled", "test-tmpl")
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", "foo", "", "", "", "").Return(nil, errors.New("random error")).Once()
		result, err := dummyService.Get("foo", "", "", "", "")
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Get", "foo", "", "", "", "")
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
