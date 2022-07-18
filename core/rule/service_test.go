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
		Rule        *rule.Rule
		Setup       func(*mocks.RuleRepository, *mocks.TemplateService, *mocks.NamespaceService, *mocks.ProviderService)
		ErrString   string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if template service return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if namespace service return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if provider service return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if list repository return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(errors.New("some error"))
					rr.EXPECT().Rollback(ctx, errors.New("some error")).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if upsert repository return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("rule.Filter")).Return(nil, errors.New("some error"))
					rr.EXPECT().Rollback(ctx, errors.New("some error")).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return nil error if upsert repository success",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{
						Variables: []template.Variable{
							{
								Name:        "var1",
								Type:        "type",
								Default:     "value",
								Description: "description",
							},
						},
					}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("rule.Filter")).Return([]rule.Rule{{
						Name:      "foo",
						Namespace: "namespace",
						Variables: []rule.RuleVariable{
							{
								Name:        "var1",
								Type:        "type",
								Value:       "value",
								Description: "description",
							},
						},
					}}, nil)
					rr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return nil error if upsert repository success and rule has no variables",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{
						Variables: []template.Variable{
							{
								Name:        "var1",
								Type:        "type",
								Default:     "value",
								Description: "description",
							},
						},
					}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("rule.Filter")).Return([]rule.Rule{{
						Name:      "foo",
						Namespace: "namespace",
						Variables: []rule.RuleVariable{
							{
								Name:        "var1",
								Type:        "type",
								Value:       "value",
								Description: "description",
							},
						},
					}}, nil)
					rr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return error if transaction commit return error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Value:       "value",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ps *mocks.ProviderService) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{
						Variables: []template.Variable{
							{
								Name:        "var1",
								Type:        "type",
								Default:     "value",
								Description: "description",
							},
						},
					}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("rule.Filter")).Return([]rule.Rule{{
						Name:      "foo",
						Namespace: "namespace",
						Variables: []rule.RuleVariable{
							{
								Name:        "var1",
								Type:        "type",
								Value:       "value",
								Description: "description",
							},
						},
					}}, nil)
					rr.EXPECT().Commit(ctx).Return(errors.New("some commit error"))
				},
				ErrString: "some commit error",
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

			err := svc.Upsert(ctx, tc.Rule)
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
		RulesWithinGroup []rule.Rule
		Setup            func(*mocks.TemplateService, *mocks.CortexClient)
		ErrString        string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if render return error",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if rendered body is empty and delete rule group return error",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("some error"))
				},
				ErrString: "error calling cortex: some error",
			},
			{
				Description: "should return nil error if rendered body is empty and delete rule group return error \"requested resource not found\"",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("requested resource not found"))
				},
			},
			{
				Description: "should return error if rendered body is empty and delete rule group success",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return("", nil)
					cc.EXPECT().DeleteRuleGroup(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
				},
			},
			{
				Description: "should return error if rendered body is failed to be unmarshalled",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(",,,---", nil)
				},
				ErrString: "cannot parse rules to alert manage rule nodes format, check your rule or template",
			},
			{
				Description: "should return error if create rule group returns error",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
					},
				},
				Setup: func(ts *mocks.TemplateService, cc *mocks.CortexClient) {
					ts.EXPECT().Render(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("map[string]string")).Return(`- alert: "InstanceDown\\nexpr:up == 0\\nfor:5m\\nlabels:{severity:page}}"`, nil)
					cc.EXPECT().CreateRuleGroup(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("rwrulefmt.RuleGroup")).Return(errors.New("some error"))
				},
				ErrString: "error calling cortex: some error",
			},
			{
				Description: "should return nil error if create rule group success",
				RulesWithinGroup: []rule.Rule{
					{
						Name:    "foo",
						Enabled: true,
						Variables: []rule.RuleVariable{
							{
								Name:        "var1",
								Type:        "type",
								Value:       "value",
								Description: "description",
							},
						},
					},
					{
						Name:    "bar",
						Enabled: false,
					},
				},
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

			err := rule.PostRuleGroupWithCortex(ctx, cortexClientMock, templateServiceMock, "namespace", "group-name", "tenant-name", tc.RulesWithinGroup)
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
