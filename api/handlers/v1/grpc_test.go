package v1

import (
	"context"
	"errors"
	"testing"

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
