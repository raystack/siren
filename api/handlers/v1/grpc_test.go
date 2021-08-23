package v1

import (
	"errors"
	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
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
		res, err := dummyGRPCServer.GetAlertHistory(nil, dummyReq)
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
		res, err := dummyGRPCServer.GetAlertHistory(nil, dummyReq)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = resource name cannot be empty")
		assert.Nil(t, res)
	})

	t.Run("should return error code 13 if getting alert history failed", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertHistoryService: mockedAlertHistoryService,
		}, logger: zaptest.NewLogger(t)}

		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).
			Return(nil, errors.New("random error")).Once()

		dummyReq := &pb.GetAlertHistoryRequest{
			Resource:  "foo",
			StartTime: 100,
			EndTime:   200,
		}
		res, err := dummyGRPCServer.GetAlertHistory(nil, dummyReq)
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
		res, err := dummyGRPCServer.GetWorkspaceChannels(nil, dummyReq)
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
		dummyGRPCServer := GRPCServer{container: &service.Container{
			WorkspaceService: mockedWorkspaceService,
		}, logger: zaptest.NewLogger(t)}

		dummyReq := &pb.GetWorkspaceChannelsRequest{
			WorkspaceName: "random",
		}
		mockedWorkspaceService.On("GetChannels", "random").
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetWorkspaceChannels(nil, dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedWorkspaceService.AssertCalled(t, "GetChannels", "random")
	})
}
