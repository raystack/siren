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
		Setup       func(*mocks.RuleRepository, *mocks.TemplateService, *mocks.NamespaceService, *mocks.RuleUploader)
		ErrString   string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, errors.New("some error"))
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{
						Provider: provider.Provider{
							Type: provider.TypeCortex,
						},
					}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(errors.New("some error"))
					rr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if upsert repository return error and rollback error",
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(errors.New("some error"))
					rr.EXPECT().Rollback(ctx, errors.New("some error")).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should commit if upsert repository and post rule return nil error",
				Rule: &rule.Rule{
					Name:      "foo",
					Namespace: "namespace",
					Variables: []rule.RuleVariable{
						{
							Name:        "var1",
							Type:        "type",
							Description: "description",
						},
					},
				},
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
					ts.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{
						Provider: provider.Provider{
							Type: provider.TypeCortex,
						},
					}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					ru.EXPECT().UpsertRule(mock.Anything, mock.AnythingOfType("*rule.Rule"), mock.AnythingOfType("*template.Template"), mock.AnythingOfType("string")).Return(nil)
					rr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return nil error if upsert repository success but post rule failed and rollback success",
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
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
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{
						Provider: provider.Provider{
							Type: provider.TypeCortex,
						},
					}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					ru.EXPECT().UpsertRule(mock.Anything, mock.AnythingOfType("*rule.Rule"), mock.AnythingOfType("*template.Template"), mock.AnythingOfType("string")).Return(errors.New("some error"))
					rr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return nil error if upsert repository success but post rule failed and rollback failed",
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
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
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{
						Provider: provider.Provider{
							Type: provider.TypeCortex,
						},
					}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					ru.EXPECT().UpsertRule(mock.Anything, mock.AnythingOfType("*rule.Rule"), mock.AnythingOfType("*template.Template"), mock.AnythingOfType("string")).Return(errors.New("some error"))
					rr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("rollback error"))
				},
				ErrString: "rollback error",
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
				Setup: func(rr *mocks.RuleRepository, ts *mocks.TemplateService, ns *mocks.NamespaceService, ru *mocks.RuleUploader) {
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
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{
						Provider: provider.Provider{
							Type: provider.TypeCortex,
						},
					}, nil)

					rr.EXPECT().WithTransaction(ctx).Return(ctx)
					rr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*rule.Rule")).Return(nil)
					ru.EXPECT().UpsertRule(mock.Anything, mock.AnythingOfType("*rule.Rule"), mock.AnythingOfType("*template.Template"), mock.AnythingOfType("string")).Return(nil)
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
				ruleUploaderMock     = new(mocks.RuleUploader)
			)
			svc := rule.NewService(
				repositoryMock,
				templateServiceMock,
				namespaceServiceMock,
				map[string]rule.RuleUploader{
					provider.TypeCortex: ruleUploaderMock,
				},
			)

			tc.Setup(repositoryMock, templateServiceMock, namespaceServiceMock, ruleUploaderMock)

			err := svc.Upsert(ctx, tc.Rule)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			templateServiceMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			ruleUploaderMock.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	t.Run("should call repository List method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.RuleRepository{}
		dummyService := rule.NewService(repositoryMock, nil, nil, nil)
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
		dummyService := rule.NewService(repositoryMock, nil, nil, nil)
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
