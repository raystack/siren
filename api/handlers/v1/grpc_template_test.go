package v1

import (
	"context"
	"errors"
	sirenv1beta1 "github.com/odpf/siren/api/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestGRPCServer_ListTemplates(t *testing.T) {
	t.Run("should return list of all templates", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.ListTemplatesRequest{}
		dummyResult := []domain.Template{
			{
				ID:   1,
				Name: "foo",
				Body: "bar",
				Tags: []string{"foo", "bar"},
				Variables: []domain.Variable{
					{
						Name:        "foo",
						Type:        "bar",
						Default:     "",
						Description: "",
					},
				},
			},
		}

		mockedTemplatesService.
			On("Index", "").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListTemplates(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetTemplates()))
		assert.Equal(t, "foo", res.GetTemplates()[0].GetName())
		assert.Equal(t, "bar", res.GetTemplates()[0].GetBody())
		assert.Equal(t, 1, len(res.GetTemplates()[0].GetVariables()))
		assert.Equal(t, "foo", res.GetTemplates()[0].GetVariables()[0].GetName())
	})

	t.Run("should return list of all templates matching particular tag", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.ListTemplatesRequest{
			Tag: "foo",
		}

		dummyResult := []domain.Template{
			{
				ID:   1,
				Name: "foo",
				Body: "bar",
				Tags: []string{"foo", "bar"},
				Variables: []domain.Variable{
					{
						Name:        "foo",
						Type:        "bar",
						Default:     "",
						Description: "",
					},
				},
			},
		}

		mockedTemplatesService.
			On("Index", "foo").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListTemplates(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetTemplates()))
		assert.Equal(t, "foo", res.GetTemplates()[0].GetName())
		assert.Equal(t, "bar", res.GetTemplates()[0].GetBody())
		assert.Equal(t, 1, len(res.GetTemplates()[0].GetVariables()))
		assert.Equal(t, "foo", res.GetTemplates()[0].GetVariables()[0].GetName())
	})

	t.Run("should return error code 13 if getting templates failed", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.ListTemplatesRequest{
			Tag: "foo",
		}
		mockedTemplatesService.
			On("Index", "foo").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListTemplates(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_GetTemplateByName(t *testing.T) {
	t.Run("should return template by name", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.GetTemplateByNameRequest{
			Name: "foo",
		}
		dummyResult := &domain.Template{
			ID:   1,
			Name: "foo",
			Body: "bar",
			Tags: []string{"foo", "bar"},
			Variables: []domain.Variable{
				{
					Name:        "foo",
					Type:        "bar",
					Default:     "",
					Description: "",
				},
			},
		}

		mockedTemplatesService.
			On("GetByName", "foo").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetTemplateByName(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), res.GetTemplate().GetId())
		assert.Equal(t, "foo", res.GetTemplate().GetName())
		assert.Equal(t, "bar", res.GetTemplate().GetBody())
		assert.Equal(t, "foo", res.GetTemplate().GetVariables()[0].GetName())
		mockedTemplatesService.AssertCalled(t, "GetByName", dummyReq.Name)
	})

	t.Run("should return error code 13 if getting template by name failed", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.GetTemplateByNameRequest{
			Name: "foo",
		}
		mockedTemplatesService.
			On("GetByName", "foo").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetTemplateByName(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_UpsertTemplate(t *testing.T) {
	dummyReq := &sirenv1beta1.UpsertTemplateRequest{
		Id:   1,
		Name: "foo",
		Body: "bar",
		Tags: []string{"foo", "bar"},
		Variables: []*sirenv1beta1.TemplateVariables{
			{
				Name:        "foo",
				Type:        "bar",
				Default:     "",
				Description: "",
			},
		},
	}
	template := &domain.Template{
		ID:   1,
		Name: "foo",
		Body: "bar",
		Tags: []string{"foo", "bar"},
		Variables: []domain.Variable{
			{
				Name:        "foo",
				Type:        "bar",
				Default:     "",
				Description: "",
			},
		},
	}

	t.Run("should return template by name", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedTemplatesService.
			On("Upsert", template).
			Return(template, nil).Once()
		res, err := dummyGRPCServer.UpsertTemplate(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), res.GetTemplate().GetId())
		assert.Equal(t, "foo", res.GetTemplate().GetName())
		assert.Equal(t, "bar", res.GetTemplate().GetBody())
		assert.Equal(t, "foo", res.GetTemplate().GetVariables()[0].GetName())
		mockedTemplatesService.AssertCalled(t, "Upsert", template)
	})

	t.Run("should return error code 13 if upsert template failed", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedTemplatesService.
			On("Upsert", template).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpsertTemplate(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_DeleteTemplate(t *testing.T) {
	t.Run("should delete template", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.DeleteTemplateRequest{
			Name: "foo",
		}

		mockedTemplatesService.
			On("Delete", "foo").
			Return(nil).Once()
		res, err := dummyGRPCServer.DeleteTemplate(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, &sirenv1beta1.DeleteTemplateResponse{}, res)
		mockedTemplatesService.AssertCalled(t, "Delete", "foo")
	})

	t.Run("should return error code 13 if deleting template failed", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1beta1.DeleteTemplateRequest{
			Name: "foo",
		}
		mockedTemplatesService.
			On("Delete", "foo").
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteTemplate(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_RenderTemplate(t *testing.T) {
	dummyReq := &sirenv1beta1.RenderTemplateRequest{
		Name: "foo",
		Variables: map[string]string{
			"foo": "bar",
		},
	}

	t.Run("should render template", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedTemplatesService.
			On("Render", "foo", dummyReq.GetVariables()).
			Return("random", nil).Once()
		res, err := dummyGRPCServer.RenderTemplate(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "random", res.GetBody())
		mockedTemplatesService.AssertCalled(t, "Render", "foo", dummyReq.GetVariables())
	})

	t.Run("should return error code 13 if rendering template failed", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedTemplatesService.
			On("Render", "foo", dummyReq.GetVariables()).
			Return("", errors.New("random error")).Once()
		res, err := dummyGRPCServer.RenderTemplate(context.Background(), dummyReq)
		assert.Empty(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}
