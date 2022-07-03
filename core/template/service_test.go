package template_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/core/template/mocks"
	"github.com/stretchr/testify/mock"
)

func TestService_Upsert(t *testing.T) {
	type testCase struct {
		Description string
		Tmpl        *template.Template
		Setup       func(*mocks.TemplateRepository)
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if upsert repository error",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), &template.Template{
					ID:   1,
					Name: "template-1",
					Body: "body of a template",
				}).Return(errors.New("some error"))
			},
			Tmpl: &template.Template{
				ID:   1,
				Name: "template-1",
				Body: "body of a template",
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error conflict if upsert repository return error duplicate",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), &template.Template{
					ID:   1,
					Name: "template-1",
					Body: "body of a template",
				}).Return(template.ErrDuplicate)
			},
			Tmpl: &template.Template{
				ID:   1,
				Name: "template-1",
				Body: "body of a template",
			},
			Err: errors.New("name already exist"),
		},
		{
			Description: "should return nil error if upsert repository not error",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().Upsert(mock.AnythingOfType("*context.emptyCtx"), &template.Template{
					ID:   1,
					Name: "template-1",
					Body: "body of a template",
				}).Return(nil)
			},
			Tmpl: &template.Template{
				ID:   1,
				Name: "template-1",
				Body: "body of a template",
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.TemplateRepository)
			)
			svc := template.NewService(repositoryMock)

			tc.Setup(repositoryMock)

			err := svc.Upsert(context.TODO(), tc.Tmpl)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestService_GetTemplate(t *testing.T) {
	var templateName = "a-template"
	type testCase struct {
		Description string
		Setup       func(*mocks.TemplateRepository)
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if get repository return error",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), templateName).Return(nil, errors.New("some error"))
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error not found if get repository return error not found",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), templateName).Return(nil, template.NotFoundError{})
			},
			Err: errors.New("template not found"),
		},
		{
			Description: "should return nil error if get by name repository not error",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), templateName).Return(&template.Template{
					ID:   1,
					Name: "template-1",
					Body: "body of a template",
				}, nil)
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.TemplateRepository)
			)
			svc := template.NewService(repositoryMock)

			tc.Setup(repositoryMock)

			_, err := svc.GetByName(context.TODO(), templateName)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}
func TestService_Render(t *testing.T) {
	type testCase struct {
		Description          string
		RequestVariables     map[string]string
		Setup                func(*mocks.TemplateRepository)
		ExpectedRenderedBody string
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should render template body from the input",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{
					Name: "foo",
					Body: "The quick [[.color]] fox jumped over the [[.adjective]] dog.",
					Tags: []string{"baz"},
					Variables: []template.Variable{
						{
							Name:        "color",
							Default:     "brown",
							Type:        "string",
							Description: "test",
						},
						{
							Name:        "adjective",
							Default:     "lazy",
							Type:        "string",
							Description: "test",
						},
					},
				}, nil)
			},
			RequestVariables: map[string]string{
				"color":     "red",
				"adjective": "dumb",
			},
			ExpectedRenderedBody: "The quick red fox jumped over the dumb dog.",
		},
		{
			Description: "should render template body enriched with defaults",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(&template.Template{
					Name: "foo",
					Body: "The quick [[.color]] fox jumped over the [[.adjective]] dog.",
					Tags: []string{"baz"},
					Variables: []template.Variable{
						{
							Name:        "color",
							Default:     "brown",
							Type:        "string",
							Description: "test",
						},
						{
							Name:        "adjective",
							Default:     "lazy",
							Type:        "string",
							Description: "test",
						},
					},
				}, nil)
			},
			RequestVariables: map[string]string{
				"adjective": "dumb",
			},
			ExpectedRenderedBody: "The quick brown fox jumped over the dumb dog.",
		},
		{
			Description: "should return not found if name does not exist",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(nil, template.NotFoundError{})
			},
			ErrString: "template not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.TemplateRepository)
			)
			svc := template.NewService(repositoryMock)

			tc.Setup(repositoryMock)

			got, err := svc.Render(context.TODO(), "foo", tc.RequestVariables)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedRenderedBody) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedRenderedBody)
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}
