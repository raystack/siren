package v1

import (
	"context"
	"errors"
	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
)


func TestGRPCServer_ListWorkspaceChannels(t *testing.T) {
	t.Run("should return workspace data object", func(t *testing.T) {
		mockedWorkspaceService := &mocks.WorkspaceService{}
		dummyResult := []domain.Channel{
			{Name: "foo", ID: "1"},
			{Name: "bar", ID: "2"},
		}
		dummyGRPCServer := GRPCServer{container: &service.Container{
			SlackWorkspaceService: mockedWorkspaceService,
		}}

		dummyReq := &sirenv1.ListWorkspaceChannelsRequest{
			WorkspaceName: "random",
		}
		mockedWorkspaceService.On("GetChannels", "random").Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListWorkspaceChannels(context.Background(), dummyReq)
		assert.Equal(t, 2, len(res.GetData()))
		assert.Equal(t, "1", res.GetData()[0].GetId())
		assert.Equal(t, "foo", res.GetData()[0].GetName())
		assert.Equal(t, "2", res.GetData()[1].GetId())
		assert.Equal(t, "bar", res.GetData()[1].GetName())
		assert.Nil(t, err)
		mockedWorkspaceService.AssertCalled(t, "GetChannels", "random")
	})

	t.Run("should return error code 13 if getting workspaces failed", func(t *testing.T) {
		mockedWorkspaceService := &mocks.WorkspaceService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				SlackWorkspaceService: mockedWorkspaceService,
			}, logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.ListWorkspaceChannelsRequest{
			WorkspaceName: "random",
		}
		mockedWorkspaceService.On("GetChannels", "random").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListWorkspaceChannels(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedWorkspaceService.AssertCalled(t, "GetChannels", "random")
	})
}

func TestGRPCServer_ExchangeCode(t *testing.T) {
	t.Run("should return OK response object", func(t *testing.T) {
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		dummyPayload := domain.OAuthPayload{Code: "foo", Workspace: "bar"}
		dummyResult := domain.OAuthExchangeResponse{OK: true}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				CodeExchangeService: mockedCodeExchangeService,
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.ExchangeCodeRequest{
			Code:      "foo",
			Workspace: "bar",
		}
		mockedCodeExchangeService.On("Exchange", dummyPayload).Return(&dummyResult, nil).Once()
		res, err := dummyGRPCServer.ExchangeCode(context.Background(), dummyReq)
		assert.Equal(t, true, res.GetOk())
		assert.Nil(t, err)
		mockedCodeExchangeService.AssertCalled(t, "Exchange", dummyPayload)
	})

	t.Run("should return error code 13 if exchange code failed", func(t *testing.T) {
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		dummyPayload := domain.OAuthPayload{Code: "foo", Workspace: "bar"}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				CodeExchangeService: mockedCodeExchangeService,
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.ExchangeCodeRequest{
			Code:      "foo",
			Workspace: "bar",
		}
		mockedCodeExchangeService.On("Exchange", dummyPayload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ExchangeCode(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedCodeExchangeService.AssertCalled(t, "Exchange", dummyPayload)
	})
}

func TestGRPCServer_GetAlertCredentials(t *testing.T) {
	t.Run("should return alert credentials of the team", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyResult := domain.AlertCredential{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "foo",
				},
				Warning: domain.SlackCredential{
					Channel: "bar",
				},
			},
		}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.GetAlertCredentialsRequest{
			TeamName: "foo",
		}
		mockedAlertmanagerService.On("Get", "foo").Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetAlertCredentials(context.Background(), dummyReq)
		assert.Equal(t, "foo", res.GetEntity())
		assert.Equal(t, "bar", res.GetTeamName())
		assert.Equal(t, "pager", res.GetPagerdutyCredentials())
		assert.Equal(t, "foo", res.GetSlackConfig().GetCritical().GetChannel())
		assert.Equal(t, "bar", res.GetSlackConfig().GetWarning().GetChannel())
		assert.Nil(t, err)
		mockedAlertmanagerService.AssertCalled(t, "Get", "foo")
	})

	t.Run("should return error code 13 if getting alert credentials failed", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.GetAlertCredentialsRequest{
			TeamName: "foo",
		}
		mockedAlertmanagerService.On("Get", "foo").
			Return(domain.AlertCredential{}, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetAlertCredentials(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedAlertmanagerService.AssertCalled(t, "Get", "foo")
	})
}

func TestGRPCServer_UpdateAlertCredentials(t *testing.T) {
	t.Run("should update alert credentials of the team", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyPayload := domain.AlertCredential{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "foo",
				},
				Warning: domain.SlackCredential{
					Channel: "bar",
				},
			},
		}
		dummyReq := &sirenv1.UpdateAlertCredentialsRequest{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &sirenv1.SlackConfig{
				Critical: &sirenv1.Critical{
					Channel: "foo",
				},
				Warning: &sirenv1.Warning{
					Channel: "bar",
				},
			},
		}
		mockedAlertmanagerService.On("Upsert", dummyPayload).Return(nil).Once()
		result, err := dummyGRPCServer.UpdateAlertCredentials(context.Background(), dummyReq)
		assert.Equal(t, result, &sirenv1.UpdateAlertCredentialsResponse{})
		assert.Nil(t, err)
		mockedAlertmanagerService.AssertCalled(t, "Upsert", dummyPayload)
	})

	t.Run("should return error code 13 if getting alert credentials failed", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyPayload := domain.AlertCredential{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "foo",
				},
				Warning: domain.SlackCredential{
					Channel: "bar",
				},
			},
		}
		dummyReq := &sirenv1.UpdateAlertCredentialsRequest{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &sirenv1.SlackConfig{
				Critical: &sirenv1.Critical{
					Channel: "foo",
				},
				Warning: &sirenv1.Warning{
					Channel: "bar",
				},
			},
		}
		mockedAlertmanagerService.On("Upsert", dummyPayload).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateAlertCredentials(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
	})

	t.Run("should return error code 3 if entity is missing", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1.UpdateAlertCredentialsRequest{
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &sirenv1.SlackConfig{
				Critical: &sirenv1.Critical{
					Channel: "foo",
				},
				Warning: &sirenv1.Warning{
					Channel: "bar",
				},
			},
		}
		res, err := dummyGRPCServer.UpdateAlertCredentials(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = entity cannot be empty")
		assert.Nil(t, res)
	})

	t.Run("should return error code 3 if pagerduty credentials is missing", func(t *testing.T) {
		mockedAlertmanagerService := &mocks.AlertmanagerService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertmanagerService: mockedAlertmanagerService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1.UpdateAlertCredentialsRequest{
			Entity:   "foo",
			TeamName: "bar",
			SlackConfig: &sirenv1.SlackConfig{
				Critical: &sirenv1.Critical{
					Channel: "foo",
				},
				Warning: &sirenv1.Warning{
					Channel: "bar",
				},
			},
		}
		res, err := dummyGRPCServer.UpdateAlertCredentials(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = pagerduty credential cannot be empty")
		assert.Nil(t, res)
	})
}

func TestGRPCServer_SendSlackNotification(t *testing.T) {
	t.Run("should return OK response", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		dummyPayload := &domain.SlackMessage{
			ReceiverName: "foo",
			ReceiverType: "channel",
			Entity:       "foo",
			Message:      "bar",
			Blocks: slack.Blocks{
				BlockSet: []slack.Block{
					&slack.SectionBlock{
						Type: slack.MBTSection,
						Text: &slack.TextBlockObject{
							Type: "mrkdwn",
							Text: "Hello",
						},
					},
				},
			},
		}
		dummyResult := &domain.SlackMessageSendResponse{
			OK: true,
		}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				NotifierServices: domain.NotifierServices{
					Slack: mockedSlackNotifierService,
				},
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.SendSlackNotificationRequest{
			Provider:     "slack",
			ReceiverName: "foo",
			ReceiverType: "channel",
			Entity:       "foo",
			Message:      "bar",
			Blocks: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"type": &structpb.Value{
							Kind: &structpb.Value_StringValue{
								StringValue: "section",
							},
						},
						"text": &structpb.Value{
							Kind: &structpb.Value_StructValue{
								StructValue: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"type": &structpb.Value{
											Kind: &structpb.Value_StringValue{
												StringValue: "mrkdwn",
											},
										},
										"text": &structpb.Value{
											Kind: &structpb.Value_StringValue{
												StringValue: "Hello",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		mockedSlackNotifierService.On("Notify", dummyPayload).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.SendSlackNotification(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, true, res.GetOk())
		mockedSlackNotifierService.AssertCalled(t, "Notify", dummyPayload)
	})

	t.Run("should return error code 13 if send slack notification failed", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		dummyPayload := &domain.SlackMessage{
			ReceiverName: "foo",
			ReceiverType: "channel",
			Entity:       "foo",
			Message:      "bar",
			Blocks:       slack.Blocks{},
		}
		dummyResult := &domain.SlackMessageSendResponse{
			OK: true,
		}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				NotifierServices: domain.NotifierServices{
					Slack: mockedSlackNotifierService,
				},
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.SendSlackNotificationRequest{
			Provider:     "slack",
			ReceiverName: "foo",
			ReceiverType: "channel",
			Entity:       "foo",
			Message:      "bar",
		}

		mockedSlackNotifierService.On("Notify", dummyPayload).
			Return(dummyResult, errors.New("random error")).Once()
		res, err := dummyGRPCServer.SendSlackNotification(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
	})

	t.Run("should return error code 3 if provider is missing", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				NotifierServices: domain.NotifierServices{
					Slack: mockedSlackNotifierService,
				},
			},
			logger: zaptest.NewLogger(t),
		}

		dummyReq := &sirenv1.SendSlackNotificationRequest{
			ReceiverName: "foo",
			ReceiverType: "channel",
			Entity:       "foo",
			Message:      "bar",
		}

		res, err := dummyGRPCServer.SendSlackNotification(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = provider not supported")
		assert.Nil(t, res)
	})
}

func TestGRPCServer_ListTemplates(t *testing.T) {
	t.Run("should return list of all templates", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &sirenv1.ListTemplatesRequest{}
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
		dummyReq := &sirenv1.ListTemplatesRequest{
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
		dummyReq := &sirenv1.ListTemplatesRequest{
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
		dummyReq := &sirenv1.GetTemplateByNameRequest{
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
		dummyReq := &sirenv1.GetTemplateByNameRequest{
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
	dummyReq := &sirenv1.UpsertTemplateRequest{
		Id:   1,
		Name: "foo",
		Body: "bar",
		Tags: []string{"foo", "bar"},
		Variables: []*sirenv1.TemplateVariables{
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
		dummyReq := &sirenv1.DeleteTemplateRequest{
			Name: "foo",
		}

		mockedTemplatesService.
			On("Delete", "foo").
			Return(nil).Once()
		res, err := dummyGRPCServer.DeleteTemplate(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, &sirenv1.DeleteTemplateResponse{}, res)
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
		dummyReq := &sirenv1.DeleteTemplateRequest{
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
	dummyReq := &sirenv1.RenderTemplateRequest{
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



