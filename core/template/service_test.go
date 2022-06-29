package template_test

import (
	"errors"
	"testing"

	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/core/template/mocks"
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
				tr.EXPECT().Upsert(&template.Template{
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
				tr.EXPECT().Upsert(&template.Template{
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
				tr.EXPECT().Upsert(&template.Template{
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

			err := svc.Upsert(tc.Tmpl)
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
				tr.EXPECT().GetByName(templateName).Return(nil, errors.New("some error"))
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error not found if get repository return error not found",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(templateName).Return(nil, template.NotFoundError{})
			},
			Err: errors.New("template not found"),
		},
		{
			Description: "should return nil error if get by name repository not error",
			Setup: func(tr *mocks.TemplateRepository) {
				tr.EXPECT().GetByName(templateName).Return(&template.Template{
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

			_, err := svc.GetByName(templateName)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}
