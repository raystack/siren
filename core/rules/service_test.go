package rules

import (
	context "context"
	"testing"
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/rules/mocks"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestService_Upsert(t *testing.T) {
	dummyTemplate := &template.Template{
		ID:        10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "tmpl",
		Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
		Tags:      []string{"baz"},
		Variables: []template.Variable{
			{
				Name:        "for",
				Default:     "10m",
				Description: "test",
				Type:        "string",
			},
			{
				Name:        "team",
				Default:     "gojek",
				Description: "test",
				Type:        "string"},
		},
	}
	rule := &domain.Rule{
		Id: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
		Variables: []domain.RuleVariable{
			{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			},
			{
				Name:        "for",
				Value:       "20m",
				Description: "test-description-for",
				Type:        "string",
			},
		},
		ProviderNamespace: 1,
	}
	dummyNamespace := &namespace.Namespace{Id: 1, Urn: "foo", Provider: 1}
	dummyProvider := &provider.Provider{Id: 1, Urn: "bar", Type: "cortex"}

	t.Run("should call repository Upsert method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		variablesMap := map[string]string{
			"for":       "20m",
			"test-name": "test-value",
		}
		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '20m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []domain.Rule{*rule}
		mockTemplateService.On("GetByName", rule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.On("GetNamespace", rule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, rule).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", rule.Namespace, rule.GroupName, "", rule.ProviderNamespace).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.On("Render", rule.Template, variablesMap).Return(expectedRender, nil).Once()
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should disable alerts", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		rule := &domain.Rule{
			Id: 1, Name: "bar", Enabled: false, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}

		mockTemplateService.On("GetByName", rule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.On("GetNamespace", rule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, rule).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", rule.Namespace, rule.GroupName, "", rule.ProviderNamespace).Return([]domain.Rule{*rule}, nil).Once()
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should handle deletion of non-existent rule group", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		rule := &domain.Rule{
			Id: 1, Name: "bar", Enabled: false, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
			Variables: []domain.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}

		mockTemplateService.On("GetByName", rule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.On("GetNamespace", rule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, rule).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", rule.Namespace, rule.GroupName, "", rule.ProviderNamespace).Return([]domain.Rule{*rule}, nil).Once()
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("requested resource not found"))
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if template service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := Service{templateService: mockTemplateService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.templateService.GetByName: random error")
		mockTemplateService.AssertExpectations(t)
	})

	t.Run("should return error if template not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := Service{templateService: mockTemplateService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(nil, nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "template not found")
		mockTemplateService.AssertExpectations(t)
	})

	t.Run("should return error if namespace service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		dummyService := Service{templateService: mockTemplateService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.namespaceService.GetNamespace: random error")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error if namespace not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		dummyService := Service{templateService: mockTemplateService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(nil, nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "namespace not found")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error if provider service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.providerService.GetProvider: random error")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if provider not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(nil, nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "provider not found")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from repository.Upsert returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Rollback: random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Upsert returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Upsert: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from not supported provider type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "not-supported-provider-type"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Rollback", ctx).Return(errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Rollback: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if provider type not supported", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "not-supported-provider-type"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "provider not supported")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from cortex client initialization returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return nil, errors.New("random error")
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Rollback", ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Rollback: random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if cortex client initialization returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return nil, errors.New("random error")
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "cortex client initialization: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from repository.Get returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Rollback: random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Get returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, namespaceService: mockNamespaceService, providerService: mockProviderService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Get: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if templateService.Render returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		dummyRulesInGroup := []domain.Rule{*rule}
		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.On("Render", mock.Anything, mock.AnythingOfType("map[string]string")).Return("", errors.New("random error")).Once()
		mockClient.On("CreateRuleGroup", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.postRuleGroupWith: s.templateService.Render: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if templateService.Render returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		dummyRulesInGroup := []domain.Rule{*rule}
		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.On("Render", mock.Anything, mock.AnythingOfType("map[string]string")).Return("", errors.New("random error")).Once()
		mockClient.On("CreateRuleGroup", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.On("Rollback", ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.repository.Rollback: random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if update cortex API call fails", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '10m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []domain.Rule{*rule}
		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.On("Render", mock.Anything, mock.AnythingOfType("map[string]string")).Return(expectedRender, nil).Once()
		mockClient.On("CreateRuleGroup", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.postRuleGroupWith: client.CreateRuleGroup: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if delete API call fails", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		mockTemplateService.On("GetByName", mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return([]domain.Rule{}, nil).Once()
		mockClient.On("DeleteRuleGroup", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.EqualError(t, err, "s.postRuleGroupWith: client.DeleteRuleGroup: random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Commit returns error", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		dummyService := Service{repository: repositoryMock, templateService: mockTemplateService, providerService: mockProviderService, namespaceService: mockNamespaceService}
		ctx := context.Background()

		mockClient := &cortexCallerMock{}
		oldCortexClientCreator := cortexClientInstance
		cortexClientInstance = func(string) (cortexCaller, error) {
			return mockClient, nil
		}
		defer func() { cortexClientInstance = oldCortexClientCreator }()

		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '20m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []domain.Rule{*rule}
		mockTemplateService.On("GetByName", mock.Anything).Return(dummyTemplate, nil).Once()
		mockNamespaceService.On("GetNamespace", mock.Anything).Return(dummyNamespace, nil).Once()
		mockProviderService.On("GetProvider", mock.Anything).Return(dummyProvider, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Upsert", ctx, mock.AnythingOfType("*domain.Rule")).Return(nil).Once()
		repositoryMock.On("Get", ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(expectedRender, nil).Once()
		mockClient.On("CreateRuleGroup", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		repositoryMock.On("Commit", ctx).Return(errors.New("random commit error")).Once()

		err := dummyService.Upsert(ctx, rule)
		assert.Error(t, err, "s.repository.Rollback: random commit error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		ctx := context.Background()

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

		repositoryMock.On("Get", ctx, "foo", "gojek", "test-group", "test-tmpl", uint64(1)).
			Return(dummyRules, nil).Once()

		result, err := dummyService.
			Get(ctx, "foo", "gojek", "test-group", "test-tmpl", 1)
		assert.Nil(t, err)
		assert.Equal(t, dummyRules, result)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &RuleRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		ctx := context.Background()

		repositoryMock.On("Get", ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, uint64(0)).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.Get(ctx, "foo", "", "", "", 0)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
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
