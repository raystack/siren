package v1

import (
	"context"
	"errors"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"strings"
	"testing"
	"time"

	sirenv1 "github.com/odpf/siren/api/proto/odpf/siren/v1"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestGRPCServer_ListAlertHistory(t *testing.T) {
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

		dummyReq := &sirenv1.ListAlertHistoryRequest{
			Resource:  "foo",
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.ListAlertHistory(context.Background(), dummyReq)
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

		dummyReq := &sirenv1.ListAlertHistoryRequest{
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.ListAlertHistory(context.Background(), dummyReq)
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

		dummyReq := &sirenv1.ListAlertHistoryRequest{
			Resource:  "foo",
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.ListAlertHistory(context.Background(), dummyReq)
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
	dummyReq := &sirenv1.CreateAlertHistoryRequest{
		Alerts: []*sirenv1.Alert{
			{
				Status: "foo",
				Labels: &sirenv1.Labels{
					Severity: "CRITICAL",
				},
				Annotations: &sirenv1.Annotations{
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

	t.Run("should not return error if parameters are missing", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertHistoryService: mockedAlertHistoryService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedAlertHistoryService.On("Create", payload).
			Return(dummyAlerts, errors.New("alert history parameters missing")).Once()

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
}

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

func TestGRPCServer_ListRules(t *testing.T) {
	dummyPayload := &sirenv1.ListRulesRequest{
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
		res, err := dummyGRPCServer.ListRules(context.Background(), dummyPayload)
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
		res, err := dummyGRPCServer.ListRules(context.Background(), dummyPayload)
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
	dummyReq := &sirenv1.UpdateRuleRequest{
		Id:        1,
		Name:      "foo",
		Namespace: "test",
		Entity:    "odpf",
		GroupName: "foo",
		Template:  "foo",
		Status:    "enabled",
		Variables: []*sirenv1.Variables{
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

		assert.Equal(t, uint64(1), res.GetRule().GetId())
		assert.Equal(t, "foo", res.GetRule().GetName())
		assert.Equal(t, "odpf", res.GetRule().GetEntity())
		assert.Equal(t, "test", res.GetRule().GetNamespace())
		assert.Equal(t, "disabled", res.GetRule().GetStatus())
		assert.Equal(t, 1, len(res.GetRule().GetVariables()))
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

func TestGRPCServer_ListProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	t.Run("should return list of all provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := []*domain.Provider{
			{
				Id:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.
			On("ListProviders").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetProviders()))
		assert.Equal(t, "foo", res.GetProviders()[0].GetHost())
		assert.Equal(t, "bar", res.GetProviders()[0].GetType())
		assert.Equal(t, "foo", res.GetProviders()[0].GetName())
	})

	t.Run("should return error code 13 if getting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("ListProviders").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := []*domain.Provider{
			{
				Id:          1,
				Host:        "foo",
				Type:        "bar",
				Name:        "foo",
				Credentials: credentials,
				Labels:      labels,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockedProviderService.
			On("ListProviders").
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListProviders(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_CreateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &domain.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1.CreateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should create provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("CreateProvider", payload).
			Return(payload, nil).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 13 if creating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("CreateProvider", payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		newPayload := &domain.Provider{
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
		}

		mockedProviderService.
			On("CreateProvider", mock.Anything).
			Return(newPayload, nil).Once()
		res, err := dummyGRPCServer.CreateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_GetProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	providerId := uint64(1)
	dummyReq := &sirenv1.GetProviderRequest{
		Id: 1,
	}

	t.Run("should return a provider", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, err)

		assert.Equal(t, "foo", res.GetHost())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 5 if no provider found", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(nil, nil).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = provider not found")
	})

	t.Run("should return error code 13 if getting provider failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, errors.New("random error")).Once()

		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		dummyResult := &domain.Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockedProviderService.
			On("GetProvider", providerId).
			Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_UpdateProvider(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	credentialsData, _ := structpb.NewStruct(credentials)

	payload := &domain.Provider{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
	}
	dummyReq := &sirenv1.UpdateProviderRequest{
		Host:        "foo",
		Type:        "bar",
		Name:        "foo",
		Credentials: credentialsData,
		Labels:      labels,
	}

	t.Run("should update provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("UpdateProvider", payload).
			Return(payload, nil).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "foo", res.GetName())
		assert.Equal(t, "bar", res.GetCredentials().GetFields()["foo"].GetStringValue())
	})

	t.Run("should return error code 13 if updating providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("UpdateProvider", payload).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		credentials["bar"] = string([]byte{0xff})
		newPayload := &domain.Provider{
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
		}

		mockedProviderService.
			On("UpdateProvider", mock.Anything).
			Return(newPayload, nil).Once()
		res, err := dummyGRPCServer.UpdateProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_DeleteProvider(t *testing.T) {
	providerId := uint64(10)
	dummyReq := &sirenv1.DeleteProviderRequest{
		Id: uint64(10),
	}

	t.Run("should delete provider object", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("DeleteProvider", providerId).
			Return(nil).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
	})

	t.Run("should return error code 13 if deleting providers failed", func(t *testing.T) {
		mockedProviderService := &mocks.ProviderService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ProviderService: mockedProviderService,
			},
			logger: zaptest.NewLogger(t),
		}

		mockedProviderService.
			On("DeleteProvider", providerId).
			Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteProvider(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}

func TestGRPCServer_ListReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"
	dummyResult := []*domain.Receiver{
		{
			Id:             1,
			Urn:            "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	t.Run("should return list of all receiver", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("ListReceivers").
			Return(dummyResult, nil).Once()

		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetReceivers()))
	})

	t.Run("should return error code 13 if getting providers failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("ListReceivers").
			Return(nil, errors.New("random error"))

		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		configurations["foo"] = string([]byte{0xff})
		dummyResult := []*domain.Receiver{
			{
				Id:             1,
				Urn:            "foo",
				Type:           "bar",
				Labels:         labels,
				Configurations: configurations,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		mockedReceiverService.
			On("ListReceivers").
			Return(dummyResult, nil)
		res, err := dummyGRPCServer.ListReceivers(context.Background(), &emptypb.Empty{})
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_CreateReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	configurationsData, _ := structpb.NewStruct(configurations)
	dummyReq := &sirenv1.CreateReceiverRequest{
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurationsData,
	}
	payload := &domain.Receiver{
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("Should create a slack receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("CreateReceiver", payload).
			Return(payload, nil).Once()

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetUrn())
		assert.Equal(t, "slack", res.GetType())
		assert.Equal(t, "bar", res.GetLabels()["foo"])
		assert.Equal(t, "foo", res.GetConfigurations().AsMap()["client_id"])
	})

	t.Run("should return error code 3 if slack client_id configuration is missing", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		slackConfigurations := make(map[string]interface{})
		slackConfigurations["client_secret"] = "foo"
		slackConfigurations["auth_code"] = "foo"

		configurationsData, _ := structpb.NewStruct(slackConfigurations)
		dummyReq := &sirenv1.CreateReceiverRequest{
			Urn:            "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurationsData,
		}

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = receiver configuration not valid")
		assert.Nil(t, res)
	})

	t.Run("should return error code 3 if slack client_secret configuration is missing", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		slackConfigurations := make(map[string]interface{})
		slackConfigurations["client_id"] = "foo"
		slackConfigurations["auth_code"] = "foo"

		configurationsData, _ := structpb.NewStruct(slackConfigurations)
		dummyReq := &sirenv1.CreateReceiverRequest{
			Urn:            "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurationsData,
		}

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = receiver configuration not valid")
		assert.Nil(t, res)
	})

	t.Run("should return error code 3 if slack auth_code configuration is missing", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		slackConfigurations := make(map[string]interface{})
		slackConfigurations["client_id"] = "foo"
		slackConfigurations["client_secret"] = "foo"

		configurationsData, _ := structpb.NewStruct(slackConfigurations)
		dummyReq := &sirenv1.CreateReceiverRequest{
			Urn:            "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurationsData,
		}

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = receiver configuration not valid")
		assert.Nil(t, res)
	})

	t.Run("should return error code 13 if creating receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("CreateReceiver", payload).
			Return(nil, errors.New("random error")).Once()

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 3 if receiver is missing", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}

		configurationsData, _ := structpb.NewStruct(configurations)
		dummyReq := &sirenv1.CreateReceiverRequest{
			Urn:            "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurationsData,
		}

		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = receiver not supported")
		assert.Nil(t, res)
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}

		configurations["workspace"] = string([]byte{0xff})
		newPayload := &domain.Receiver{
			Urn:            "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
		}

		mockedReceiverService.
			On("CreateReceiver", mock.Anything).
			Return(newPayload, nil)
		res, err := dummyGRPCServer.CreateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_GetReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	receiverId := uint64(1)
	dummyReq := &sirenv1.GetReceiverRequest{
		Id: 1,
	}
	payload := &domain.Receiver{
		Urn:            "foo",
		Type:           "bar",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("should return a receiver", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("GetReceiver", receiverId).
			Return(payload, nil).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetUrn())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "bar", res.GetLabels()["foo"])
		assert.Equal(t, "bar", res.GetConfigurations().AsMap()["foo"])
	})

	t.Run("should return error code 5 if no receiver found", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("GetReceiver", receiverId).
			Return(nil, nil).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = receiver not found")
	})

	t.Run("should return error code 13 if getting receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("GetReceiver", receiverId).
			Return(payload, errors.New("random error")).Once()

		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}

		configurations["foo"] = string([]byte{0xff})
		payload := &domain.Receiver{
			Urn:            "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
		}

		mockedReceiverService.
			On("GetReceiver", receiverId).
			Return(payload, nil)
		res, err := dummyGRPCServer.GetReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_UpdateReceiver(t *testing.T) {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	configurationsData, _ := structpb.NewStruct(configurations)
	dummyReq := &sirenv1.UpdateReceiverRequest{
		Urn:            "foo",
		Type:           "bar",
		Labels:         labels,
		Configurations: configurationsData,
	}
	payload := &domain.Receiver{
		Urn:            "foo",
		Type:           "bar",
		Labels:         labels,
		Configurations: configurations,
	}

	t.Run("should update receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("UpdateReceiver", payload).
			Return(payload, nil).Once()

		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "foo", res.GetUrn())
		assert.Equal(t, "bar", res.GetType())
		assert.Equal(t, "bar", res.GetLabels()["foo"])
		assert.Equal(t, "bar", res.GetConfigurations().AsMap()["foo"])
	})

	t.Run("should return error code 13 if updating receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("UpdateReceiver", payload).
			Return(nil, errors.New("random error"))

		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})

	t.Run("should return error code 13 if NewStruct conversion failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		configurations["foo"] = string([]byte{0xff})
		newPayload := &domain.Receiver{
			Urn:            "foo",
			Type:           "bar",
			Labels:         labels,
			Configurations: configurations,
		}

		mockedReceiverService.
			On("UpdateReceiver", mock.Anything).
			Return(newPayload, nil)
		res, err := dummyGRPCServer.UpdateReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.Equal(t, strings.Replace(err.Error(), "\u00a0", " ", -1),
			"rpc error: code = Internal desc = proto: invalid UTF-8 in string: \"\\xff\"")
	})
}

func TestGRPCServer_DeleteReceiver(t *testing.T) {
	providerId := uint64(10)
	dummyReq := &sirenv1.DeleteReceiverRequest{
		Id: uint64(10),
	}

	t.Run("should delete receiver object", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("DeleteReceiver", providerId).
			Return(nil).Once()

		res, err := dummyGRPCServer.DeleteReceiver(context.Background(), dummyReq)
		assert.Nil(t, err)
		assert.Equal(t, "", res.String())
	})

	t.Run("should return error code 13 if deleting receiver failed", func(t *testing.T) {
		mockedReceiverService := &mocks.ReceiverService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				ReceiverService: mockedReceiverService,
			},
			logger: zaptest.NewLogger(t),
		}
		mockedReceiverService.
			On("DeleteReceiver", providerId).
			Return(errors.New("random error")).Once()

		res, err := dummyGRPCServer.DeleteReceiver(context.Background(), dummyReq)
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
	})
}
