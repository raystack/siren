package v1

import (
	"context"
	"errors"
	"github.com/slack-go/slack"
	"google.golang.org/protobuf/types/known/structpb"
	"testing"
	"time"

	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestGRPCServer_GetAlertHistory(t *testing.T) {
	t.Run("should return alert history objects", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}
		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).
			Return(dummyAlerts, nil).Once()
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertHistoryService: mockedAlertHistoryService,
		}}

		dummyReq := &pb.GetAlertHistoryRequest{
			Resource:  "foo",
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.GetAlertHistory(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetName())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetTemplateId())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetLevel())
		assert.Nil(t, err)
		mockedAlertHistoryService.AssertCalled(t, "Get", "foo", uint32(100), uint32(200))
	})

	t.Run("should return error code 3 if resource query param is missing", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertHistoryService: mockedAlertHistoryService,
		}}

		dummyReq := &pb.GetAlertHistoryRequest{
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.GetAlertHistory(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = resource name cannot be empty")
		assert.Nil(t, res)
	})

	t.Run("should return error code 13 if getting alert history failed", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertHistoryService: mockedAlertHistoryService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).
			Return(nil, errors.New("random error")).Once()

		dummyReq := &pb.GetAlertHistoryRequest{
			Resource:  "foo",
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.GetAlertHistory(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedAlertHistoryService.AssertCalled(t, "Get", "foo", uint32(100), uint32(200))
	})
}

func TestGRPCServer_CreateAlertHistory(t *testing.T) {
	payload := &domain.Alerts{
		Alerts: []domain.Alert{
			{
				Status: "foo",
				Labels: domain.Labels{
					Severity: "CRITICAL",
				},
				Annotations: domain.Annotations{
					Resource:    "foo",
					Template:    "random",
					MetricName:  "bar",
					MetricValue: "30",
				},
			},
		},
	}
	dummyReq := &pb.CreateAlertHistoryRequest{
		Alerts: []*pb.Alerts{
			{
				Status: "foo",
				Labels: &pb.Labels{
					Severity: "CRITICAL",
				},
				Annotations: &pb.Annotations{
					Resource:    "foo",
					Template:    "random",
					MetricName:  "bar",
					MetricValue: "30",
				},
			},
		},
	}

	t.Run("should create alert history objects", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}
		mockedAlertHistoryService.On("Create", payload).
			Return(dummyAlerts, nil).Once()
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertHistoryService: mockedAlertHistoryService,
		}}

		res, err := dummyGRPCServer.CreateAlertHistory(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetName())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetTemplateId())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetLevel())
		assert.Nil(t, err)
		mockedAlertHistoryService.AssertCalled(t, "Create", payload)
	})

	t.Run("should return error code 13 if getting alert history failed", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertHistoryService: mockedAlertHistoryService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedAlertHistoryService.On("Create", payload).
			Return(nil, errors.New("random error")).Once()

		res, err := dummyGRPCServer.CreateAlertHistory(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedAlertHistoryService.AssertCalled(t, "Create", payload)
	})

	t.Run("should return error code 3 if parameters is missing", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertHistoryService: mockedAlertHistoryService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedAlertHistoryService.On("Create", payload).
			Return(nil, errors.New("alert history parameters missing")).Once()

		res, err := dummyGRPCServer.CreateAlertHistory(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = alert history parameters missing")
		assert.Nil(t, res)
		mockedAlertHistoryService.AssertCalled(t, "Create", payload)
	})
}

func TestGRPCServer_GetWorkspaceChannels(t *testing.T) {
	t.Run("should return workspace data object", func(t *testing.T) {
		mockedWorkspaceService := &mocks.WorkspaceService{}
		dummyResult := []domain.Channel{
			{Name: "foo", ID: "1"},
			{Name: "bar", ID: "2"},
		}
		dummyGRPCServer := GRPCServer{container: &service.Container{
			WorkspaceService: mockedWorkspaceService,
		}}

		dummyReq := &pb.GetWorkspaceChannelsRequest{
			WorkspaceName: "random",
		}
		mockedWorkspaceService.On("GetChannels", "random").Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetWorkspaceChannels(context.Background(), dummyReq)
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
				WorkspaceService: mockedWorkspaceService,
			}, logger: zaptest.NewLogger(t),
		}

		dummyReq := &pb.GetWorkspaceChannelsRequest{
			WorkspaceName: "random",
		}
		mockedWorkspaceService.On("GetChannels", "random").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetWorkspaceChannels(context.Background(), dummyReq)
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

		dummyReq := &pb.ExchangeCodeRequest{
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

		dummyReq := &pb.ExchangeCodeRequest{
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

		dummyReq := &pb.GetAlertCredentialsRequest{
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

		dummyReq := &pb.GetAlertCredentialsRequest{
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
		dummyReq := &pb.UpdateAlertCredentialsRequest{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &pb.SlackConfig{
				Critical: &pb.Critical{
					Channel: "foo",
				},
				Warning: &pb.Warning{
					Channel: "bar",
				},
			},
		}
		mockedAlertmanagerService.On("Upsert", dummyPayload).Return(nil).Once()
		result, err := dummyGRPCServer.UpdateAlertCredentials(context.Background(), dummyReq)
		assert.Equal(t, result, &pb.UpdateAlertCredentialsResponse{})
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
		dummyReq := &pb.UpdateAlertCredentialsRequest{
			Entity:               "foo",
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &pb.SlackConfig{
				Critical: &pb.Critical{
					Channel: "foo",
				},
				Warning: &pb.Warning{
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
		dummyReq := &pb.UpdateAlertCredentialsRequest{
			TeamName:             "bar",
			PagerdutyCredentials: "pager",
			SlackConfig: &pb.SlackConfig{
				Critical: &pb.Critical{
					Channel: "foo",
				},
				Warning: &pb.Warning{
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
		dummyReq := &pb.UpdateAlertCredentialsRequest{
			Entity:   "foo",
			TeamName: "bar",
			SlackConfig: &pb.SlackConfig{
				Critical: &pb.Critical{
					Channel: "foo",
				},
				Warning: &pb.Warning{
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

		dummyReq := &pb.SendSlackNotificationRequest{
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

		dummyReq := &pb.SendSlackNotificationRequest{
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

		dummyReq := &pb.SendSlackNotificationRequest{
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

func TestGRPCServer_GetRules(t *testing.T) {
	dummyPayload := &pb.GetRulesRequest{
		Namespace: "test",
		Entity:    "odpf",
		GroupName: "foo",
		Status:    "enabled",
		Template:  "foo",
	}

	t.Run("should return stored rules", func(t *testing.T) {
		mockedRuleService := &mocks.RuleService{}
		dummyResult := []domain.Rule{
			{
				ID:        1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      "foo",
				Namespace: "test",
				Entity:    "odpf",
				GroupName: "foo",
				Template:  "foo",
				Status:    "enabled",
				Variables: []domain.RuleVariable{
					{
						Name:        "foo",
						Type:        "int",
						Value:       "bar",
						Description: "",
					},
				},
			},
		}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				RulesService: mockedRuleService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedRuleService.
			On("Get", dummyPayload.Namespace, dummyPayload.Entity, dummyPayload.GroupName, dummyPayload.Status, dummyPayload.Template).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetRules(context.Background(), dummyPayload)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetRules()))
		assert.Equal(t, uint64(1), res.GetRules()[0].GetId())
		assert.Equal(t, "foo", res.GetRules()[0].GetName())
		assert.Equal(t, "odpf", res.GetRules()[0].GetEntity())
		assert.Equal(t, "test", res.GetRules()[0].GetNamespace())
		assert.Equal(t, "enabled", res.GetRules()[0].GetStatus())
		assert.Equal(t, 1, len(res.GetRules()[0].GetVariables()))
		mockedRuleService.AssertCalled(t, "Get", "test", "odpf", "foo", "enabled", "foo")
	})
	t.Run("should return error code 13 if getting rules failed", func(t *testing.T) {
		mockedRuleService := &mocks.RuleService{}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				RulesService: mockedRuleService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedRuleService.
			On("Get", dummyPayload.Namespace, dummyPayload.Entity, dummyPayload.GroupName, dummyPayload.Status, dummyPayload.Template).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetRules(context.Background(), dummyPayload)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_UpdateRules(t *testing.T) {
	dummyPayload := domain.Rule{
		ID:        1,
		Name:      "foo",
		Namespace: "test",
		Entity:    "odpf",
		GroupName: "foo",
		Template:  "foo",
		Status:    "enabled",
		Variables: []domain.RuleVariable{
			{
				Name:        "foo",
				Type:        "int",
				Value:       "bar",
				Description: "",
			},
		},
	}
	dummyReq := &pb.UpdateRuleRequest{
		Id:        1,
		Name:      "foo",
		Namespace: "test",
		Entity:    "odpf",
		GroupName: "foo",
		Template:  "foo",
		Status:    "enabled",
		Variables: []*pb.Variables{
			{
				Name:        "foo",
				Type:        "int",
				Value:       "bar",
				Description: "",
			},
		},
	}
	t.Run("should update rule", func(t *testing.T) {
		mockedRuleService := &mocks.RuleService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				RulesService: mockedRuleService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := dummyPayload
		dummyResult.Status = "disabled"

		mockedRuleService.
			On("Upsert", &dummyPayload).
			Return(&dummyResult, nil).Once()
		res, err := dummyGRPCServer.UpdateRule(context.Background(), dummyReq)
		assert.Nil(t, err)

		assert.Equal(t, uint64(1), res.GetId())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "odpf", res.GetEntity())
		assert.Equal(t, "test", res.GetNamespace())
		assert.Equal(t, "disabled", res.GetStatus())
		assert.Equal(t, 1, len(res.GetVariables()))
		mockedRuleService.AssertCalled(t, "Upsert", &dummyPayload)
	})

	t.Run("should return error code 13 if getting rules failed", func(t *testing.T) {
		mockedRuleService := &mocks.RuleService{}

		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				RulesService: mockedRuleService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedRuleService.
			On("Upsert", &dummyPayload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateRule(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_GetTemplates(t *testing.T) {
	t.Run("should return list of all templates", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				TemplatesService: mockedTemplatesService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyReq := &pb.GetTemplatesRequest{}
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
		res, err := dummyGRPCServer.GetTemplates(context.Background(), dummyReq)
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
		dummyReq := &pb.GetTemplatesRequest{
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
		res, err := dummyGRPCServer.GetTemplates(context.Background(), dummyReq)
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
		dummyReq := &pb.GetTemplatesRequest{
			Tag: "foo",
		}
		mockedTemplatesService.
			On("Index", "foo").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetTemplates(context.Background(), dummyReq)
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
		dummyReq := &pb.GetTemplateByNameRequest{
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
		assert.Equal(t, uint64(1), res.GetId())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetBody())
		assert.Equal(t, "foo", res.GetVariables()[0].GetName())
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
		dummyReq := &pb.GetTemplateByNameRequest{
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
	dummyReq := &pb.UpsertTemplateRequest{
		Id:   1,
		Name: "foo",
		Body: "bar",
		Tags: []string{"foo", "bar"},
		Variables: []*pb.TemplateVariables{
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
		assert.Equal(t, uint64(1), res.GetId())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetBody())
		assert.Equal(t, "foo", res.GetVariables()[0].GetName())
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
		dummyReq := &pb.DeleteTemplateRequest{
			Name: "foo",
		}

		mockedTemplatesService.
			On("Delete", "foo").
			Return(nil).Once()
		res, err := dummyGRPCServer.DeleteTemplate(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, &pb.DeleteTemplateResponse{}, res)
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
		dummyReq := &pb.DeleteTemplateRequest{
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
	dummyReq := &pb.RenderTemplateRequest{
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
