package rule_test

import (
	context "context"
	"testing"
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/rule/mocks"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/pkg/errors"
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
	theRule := &rule.Rule{
		ID: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
		Variables: []rule.RuleVariable{
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
	dummyNamespace := &namespace.Namespace{ID: 1, URN: "foo", Provider: 1}
	dummyProvider := &provider.Provider{ID: 1, URN: "bar", Type: "cortex"}

	t.Run("should call repository Upsert method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		variablesMap := map[string]string{
			"for":       "20m",
			"test-name": "test-value",
		}
		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '20m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []rule.Rule{*theRule}
		mockTemplateService.EXPECT().GetByName(theRule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(theRule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, theRule).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", theRule.Namespace, theRule.GroupName, "", theRule.ProviderNamespace).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.EXPECT().Render(theRule.Template, variablesMap).Return(expectedRender, nil).Once()
		mockClient.EXPECT().CreateRuleGroup(mock.Anything, "foo", mock.Anything).Return(nil)
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should disable alerts", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}

		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		theRule := &rule.Rule{
			ID: 1, Name: "bar", Enabled: false, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
			Variables: []rule.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}

		mockTemplateService.EXPECT().GetByName(theRule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(theRule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, theRule).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", theRule.Namespace, theRule.GroupName, "", theRule.ProviderNamespace).Return([]rule.Rule{*theRule}, nil).Once()
		mockClient.EXPECT().DeleteRuleGroup(mock.Anything, "foo", mock.Anything).Return(nil)
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should handle deletion of non-existent rule group", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		theRule := &rule.Rule{
			ID: 1, Name: "bar", Enabled: false, GroupName: "test-group", Namespace: "foo", Template: "test-tmpl",
			Variables: []rule.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}

		mockTemplateService.EXPECT().GetByName(theRule.Template).Return(dummyTemplate, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(theRule.ProviderNamespace).Return(dummyNamespace, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, theRule).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", theRule.Namespace, theRule.GroupName, "", theRule.ProviderNamespace).Return([]rule.Rule{*theRule}, nil).Once()
		mockClient.EXPECT().DeleteRuleGroup(mock.Anything, "foo", mock.Anything).Return(errors.New("requested resource not found"))
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if template service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := rule.NewService(nil, mockTemplateService, nil, nil, nil)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		mockTemplateService.AssertExpectations(t)
	})

	t.Run("should return error if template not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		dummyService := rule.NewService(nil, mockTemplateService, nil, nil, nil)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(nil, errors.ErrNotFound.WithMsgf("template not found")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "template not found")
		mockTemplateService.AssertExpectations(t)
	})

	t.Run("should return error if namespace service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		dummyService := rule.NewService(nil, mockTemplateService, mockNamespaceService, nil, nil)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error if namespace not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(nil, mockTemplateService, mockNamespaceService, nil, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(nil, errors.ErrNotFound.WithMsgf("namespace not found")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "namespace not found")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
	})

	t.Run("should return error if provider service returns error", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(nil, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if provider not found", func(t *testing.T) {
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(nil, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, errors.ErrNotFound.WithMsgf("provider not found")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "provider not found")
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from repository.Upsert returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Upsert returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error conflict if repository.Upsert returns err duplicate", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(rule.ErrDuplicate).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "rule conflicted with existing")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from not supported provider type", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "not-supported-provider-type"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(errors.New("random error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if provider type not supported", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "not-supported-provider-type"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "provider not supported")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if rollback from repository.Get returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Get returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if templateService.Render returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		dummyRulesInGroup := []rule.Rule{*theRule}
		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.EXPECT().Render(mock.Anything, mock.AnythingOfType("map[string]string")).Return("", errors.New("random error")).Once()
		mockClient.EXPECT().CreateRuleGroup(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if templateService.Render returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		dummyRulesInGroup := []rule.Rule{*theRule}
		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.EXPECT().Render(mock.Anything, mock.AnythingOfType("map[string]string")).Return("", errors.New("random error")).Once()
		mockClient.EXPECT().CreateRuleGroup(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.EXPECT().Rollback(ctx).Return(errors.New("random rollback error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random rollback error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if update cortex API call fails", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '10m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []rule.Rule{*theRule}
		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.EXPECT().Render(mock.Anything, mock.AnythingOfType("map[string]string")).Return(expectedRender, nil).Once()
		mockClient.EXPECT().CreateRuleGroup(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should rollback if delete API call fails", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(&template.Template{}, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(&namespace.Namespace{}, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(&provider.Provider{Type: "cortex"}, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return([]rule.Rule{}, nil).Once()
		mockClient.EXPECT().DeleteRuleGroup(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("random error"))
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})

	t.Run("should return error if repository.Commit returns error", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		mockTemplateService := &mocks.TemplatesService{}
		mockNamespaceService := &mocks.NamespaceService{}
		mockProviderService := &mocks.ProviderService{}
		mockClient := &mocks.CortexClient{}
		dummyService := rule.NewService(repositoryMock, mockTemplateService, mockNamespaceService, mockProviderService, mockClient)
		ctx := context.Background()

		expectedRender := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '20m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

		dummyRulesInGroup := []rule.Rule{*theRule}
		mockTemplateService.EXPECT().GetByName(mock.Anything).Return(dummyTemplate, nil).Once()
		mockNamespaceService.EXPECT().GetNamespace(mock.Anything).Return(dummyNamespace, nil).Once()
		mockProviderService.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(dummyProvider, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Upsert(ctx, mock.AnythingOfType("*rule.Rule")).Return(nil).Once()
		repositoryMock.EXPECT().Get(ctx, "", mock.Anything, mock.Anything, "", mock.Anything).Return(dummyRulesInGroup, nil).Once()
		mockTemplateService.EXPECT().Render(mock.Anything, mock.Anything).Return(expectedRender, nil).Once()
		mockClient.EXPECT().CreateRuleGroup(mock.Anything, mock.Anything, mock.Anything).Return(nil)
		repositoryMock.EXPECT().Commit(ctx).Return(errors.New("random commit error")).Once()

		err := dummyService.Upsert(ctx, theRule)
		assert.Error(t, err, "s.repository.Rollback: random commit error")
		repositoryMock.AssertExpectations(t)
		mockTemplateService.AssertExpectations(t)
		mockNamespaceService.AssertExpectations(t)
		mockProviderService.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		dummyService := rule.NewService(repositoryMock, nil, nil, nil, nil)
		ctx := context.Background()

		dummyRules := []rule.Rule{{
			ID: 1, Name: "bar", Enabled: true, GroupName: "test-group", Namespace: "baz", Template: "test-tmpl",
			Variables: []rule.RuleVariable{{
				Name:        "test-name",
				Value:       "test-value",
				Description: "test-description",
				Type:        "test-type",
			}},
			ProviderNamespace: 1,
		}}

		repositoryMock.EXPECT().Get(ctx, "foo", "gojek", "test-group", "test-tmpl", uint64(1)).
			Return(dummyRules, nil).Once()

		result, err := dummyService.
			Get(ctx, "foo", "gojek", "test-group", "test-tmpl", 1)
		assert.Nil(t, err)
		assert.Equal(t, dummyRules, result)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		dummyService := rule.NewService(repositoryMock, nil, nil, nil, nil)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything, uint64(0)).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.Get(ctx, "foo", "", "", "", 0)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
	})
}
