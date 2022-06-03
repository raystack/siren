package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/domain"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGRPCServer_ListAlerts(t *testing.T) {
	t.Run("should return alert history objects", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		timenow := time.Now()
		dummyAlerts := []domain.Alert{{
			Id: 1, ProviderId: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "bar", MetricValue: "30", Rule: "bar",
			TriggeredAt: timenow,
		}}
		mockedAlertService.On("Get", "foo", uint64(1), uint64(100), uint64(200)).
			Return(dummyAlerts, nil).Once()
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertService: mockedAlertService,
		}}

		dummyReq := &sirenv1beta1.ListAlertsRequest{
			ResourceName: "foo",
			ProviderId:   1,
			StartTime:    100,
			EndTime:      200,
		}
		res, err := dummyGRPCServer.ListAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetSeverity())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetRule())
		assert.Nil(t, err)
		mockedAlertService.AssertCalled(t, "Get", "foo", uint64(1), uint64(100), uint64(200))
	})

	t.Run("should return error code 13 if getting alert history failed", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertService: mockedAlertService,
			},
			logger: log.NewNoop(),
		}

		mockedAlertService.On("Get", "foo", uint64(1), uint64(100), uint64(200)).
			Return(nil, errors.New("random error")).Once()

		dummyReq := &sirenv1beta1.ListAlertsRequest{
			ResourceName: "foo",
			ProviderId:   1,
			StartTime:    100,
			EndTime:      200,
		}
		res, err := dummyGRPCServer.ListAlerts(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedAlertService.AssertCalled(t, "Get", "foo", uint64(1), uint64(100), uint64(200))
	})
}

func TestGRPCServer_CreateAlertHistory(t *testing.T) {
	timenow := timestamppb.New(time.Now())
	payload := &domain.Alerts{
		Alerts: []domain.Alert{
			{
				ProviderId:   1,
				ResourceName: "foo",
				MetricName:   "bar",
				MetricValue:  "30",
				Severity:     "CRITICAL",
				Rule:         "random",
				TriggeredAt:  timenow.AsTime(),
			},
		},
	}
	dummyReq := &sirenv1beta1.CreateCortexAlertsRequest{
		ProviderId: 1,
		Alerts: []*sirenv1beta1.CortexAlert{
			{
				Status: "foo",
				Labels: &sirenv1beta1.Labels{
					Severity: "CRITICAL",
				},
				Annotations: &sirenv1beta1.Annotations{
					Resource:    "foo",
					Template:    "random",
					MetricName:  "bar",
					MetricValue: "30",
				},
				StartsAt: timenow,
			},
		},
	}

	t.Run("should create alerts objects", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyAlerts := []domain.Alert{{
			Id:           1,
			ProviderId:   1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Severity:     "CRITICAL",
			Rule:         "random",
			TriggeredAt:  timenow.AsTime(),
		}}
		mockedAlertService.On("Create", payload).
			Return(dummyAlerts, nil).Once()
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertService: mockedAlertService,
		}}

		res, err := dummyGRPCServer.CreateCortexAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertCalled(t, "Create", payload)
	})

	t.Run("should create alerts for resolved alerts", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyReq := &sirenv1beta1.CreateCortexAlertsRequest{
			ProviderId: 1,
			Alerts: []*sirenv1beta1.CortexAlert{
				{
					Status: "resolved",
					Labels: &sirenv1beta1.Labels{
						Severity: "CRITICAL",
					},
					Annotations: &sirenv1beta1.Annotations{
						Resource:    "foo",
						Template:    "random",
						MetricName:  "bar",
						MetricValue: "30",
					},
					StartsAt: timenow,
				},
			},
		}
		payload := &domain.Alerts{
			Alerts: []domain.Alert{
				{
					ProviderId:   1,
					ResourceName: "foo",
					MetricName:   "bar",
					MetricValue:  "30",
					Severity:     "resolved",
					Rule:         "random",
					TriggeredAt:  timenow.AsTime(),
				},
			},
		}
		dummyAlerts := []domain.Alert{{
			Id:           1,
			ProviderId:   1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Severity:     "resolved",
			Rule:         "random",
			TriggeredAt:  timenow.AsTime(),
		}}
		mockedAlertService.On("Create", payload).
			Return(dummyAlerts, nil).Once()
		dummyGRPCServer := GRPCServer{container: &service.Container{
			AlertService: mockedAlertService,
		}}

		res, err := dummyGRPCServer.CreateCortexAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "resolved", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertCalled(t, "Create", payload)
	})

	t.Run("should return error code 13 if getting alert history failed", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertService: mockedAlertService,
			},
			logger: log.NewNoop(),
		}

		mockedAlertService.On("Create", payload).
			Return(nil, errors.New("random error")).Once()

		res, err := dummyGRPCServer.CreateCortexAlerts(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = random error")
		assert.Nil(t, res)
		mockedAlertService.AssertCalled(t, "Create", payload)
	})

	t.Run("should insert valid alerts and should not return error if parameters are missing", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyReq := &sirenv1beta1.CreateCortexAlertsRequest{
			ProviderId: 1,
			Alerts: []*sirenv1beta1.CortexAlert{
				&sirenv1beta1.CortexAlert{
					Status: "foo",
					Labels: &sirenv1beta1.Labels{
						Severity: "CRITICAL",
					},
					Annotations: &sirenv1beta1.Annotations{
						Resource:    "foo",
						MetricName:  "bar",
						MetricValue: "30",
					},
					StartsAt: timenow,
				},
				&sirenv1beta1.CortexAlert{
					Status: "foo",
					Labels: &sirenv1beta1.Labels{
						Severity: "CRITICAL",
					},
					Annotations: &sirenv1beta1.Annotations{
						Resource:    "foo",
						Template:    "random",
						MetricName:  "bar",
						MetricValue: "30",
					},
					StartsAt: timenow,
				},
			},
		}
		dummyAlerts := []domain.Alert{{
			ProviderId:   1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Rule:         "random",
			Severity:     "CRITICAL",
			TriggeredAt:  time.Now(),
		}}
		dummyGRPCServer := GRPCServer{
			container: &service.Container{
				AlertService: mockedAlertService,
			},
			logger: log.NewNoop(),
		}

		mockedAlertService.On("Create", payload).
			Return(dummyAlerts, nil).Once()

		res, err := dummyGRPCServer.CreateCortexAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetProviderId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertCalled(t, "Create", payload)
	})
}
