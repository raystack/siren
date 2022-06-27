package rule_test

import (
	context "context"
	"testing"

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
	type testCase struct {
		Description string
		Setup       func(*mocks.RuleRepository, *mocks.TemplateService, *mocks.NamespaceService, *mocks.ProviderService)
		ErrString   string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if template service return error",
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if namespace service return error",
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if provider service return error",
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if upsert repository return error",
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)
					rr.EXPECT().UpsertWithTx(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule"), mock.Anything).Return(0, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return nil error if upsert repository success",
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)
					rr.EXPECT().UpsertWithTx(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule"), mock.Anything).Return(1, nil)
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock       = new(mocks.RuleRepository)
				templateServiceMock  = new(mocks.TemplateService)
				namespaceServiceMock = new(mocks.NamespaceService)
				providerServiceMock  = new(mocks.ProviderService)
			)
			svc := rule.NewService(repositoryMock, templateServiceMock, namespaceServiceMock, providerServiceMock, nil)

			tc.Setup(repositoryMock, templateServiceMock, namespaceServiceMock, providerServiceMock)

			_, err := svc.Upsert(ctx, &rule.Rule{
				Name:      "foo",
				Namespace: "namespace",
			})
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			templateServiceMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			providerServiceMock.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Run("should call repository List method and return result in domain's type", func(t *testing.T) {
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

		repositoryMock.EXPECT().List(ctx, rule.Filter{
			Name:         "foo",
			Namespace:    "odpf",
			GroupName:    "test-group",
			TemplateName: "test-tmpl",
			NamespaceID:  1,
		}).Return(dummyRules, nil).Once()

		result, err := dummyService.List(ctx, rule.Filter{
			Name:         "foo",
			Namespace:    "odpf",
			GroupName:    "test-group",
			TemplateName: "test-tmpl",
			NamespaceID:  1,
		})
		assert.Nil(t, err)
		assert.Equal(t, dummyRules, result)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository List method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		dummyService := rule.NewService(repositoryMock, nil, nil, nil, nil)
		ctx := context.Background()

		repositoryMock.EXPECT().List(ctx, mock.Anything).Return(nil, errors.New("random error")).Once()

		result, err := dummyService.List(ctx, rule.Filter{
			Name: "foo",
		})
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
	})
}

func TestService_PostRuleGroupWithCortex(t *testing.T) {
	type testCase struct {
		Description      string
		Rule             *rule.Rule
		RulesWithinGroup []rule.Rule
		Setup            func(*mocks.TemplateService, *mocks.CortexClient)
		ErrString        string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if render return error",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if rendered body is empty and delete rule group return error",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("some error"))
				},
				ErrString: "error calling cortex: some error",
			},
			{
				Description: "should return nil error if rendered body is empty and delete rule group return error \"requested resource not found\"",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("requested resource not found"))
				},
			},
			{
				Description: "should return error if rendered body is empty and delete rule group success",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
				},
			},
			{
				Description: "should return error if rendered body is failed to be unmarshalled",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(",,,---", nil)
				},
				ErrString: "cannot parse rules to alert manage rule nodes format, check your rule or template",
			},
			{
				Description: "should return error if create rule group returns error",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(`- alert: "InstanceDown\\nexpr:up == 0\\nfor:5m\\nlabels:{severity:page}}"`, nil)
					cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(errors.New("some error"))
				},
				ErrString: "error calling cortex: some error",
			},
			{
				Description: "should return nil error if create rule group success",
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(`- alert: "InstanceDown\\nexpr:up == 0\\nfor:5m\\nlabels:{severity:page}}"`, nil)
					cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(nil)
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				templateServiceMock = new(mocks.TemplateService)
				cortexClientMock    = new(mocks.CortexClient)
			)

			tc.Setup(templateServiceMock, cortexClientMock)

			err := rule.PostRuleGroupWithCortex(ctx, cortexClientMock, templateServiceMock, &rule.Rule{}, []rule.Rule{
				{
					Name:    "foo",
					Enabled: true,
				},
			}, "tenant-name")
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			templateServiceMock.AssertExpectations(t)
			cortexClientMock.AssertExpectations(t)
		})
	}
}
