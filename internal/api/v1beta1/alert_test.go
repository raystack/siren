package v1beta1_test

import (
	"context"
	"testing"
	"time"

	"github.com/raystack/salt/log"
	"github.com/raystack/siren/core/alert"
	"github.com/raystack/siren/core/provider"
	"github.com/raystack/siren/internal/api"
	"github.com/raystack/siren/internal/api/mocks"
	"github.com/raystack/siren/internal/api/v1beta1"
	"github.com/raystack/siren/pkg/errors"
	sirenv1beta1 "github.com/raystack/siren/proto/raystack/siren/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestGRPCServer_ListAlerts(t *testing.T) {
	t.Run("should return alert history objects", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		timenow := time.Now()
		dummyAlerts := []alert.Alert{{
			ID: 1, ProviderID: 1, ResourceName: "foo", Severity: "CRITICAL", MetricName: "bar", MetricValue: "30", Rule: "bar",
			TriggeredAt: timenow,
		}}
		mockedAlertService.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    100,
			EndTime:      200,
		}).Return(dummyAlerts, nil).Once()
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, nil, api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService})

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
		mockedAlertService.AssertExpectations(t)
	})

	t.Run("should return error Internal if getting alert history failed", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService})

		mockedAlertService.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), alert.Filter{
			ProviderID:   1,
			ResourceName: "foo",
			StartTime:    100,
			EndTime:      200,
		}).Return(nil, errors.New("random error")).Once()

		dummyReq := &sirenv1beta1.ListAlertsRequest{
			ResourceName: "foo",
			ProviderId:   1,
			StartTime:    100,
			EndTime:      200,
		}
		res, err := dummyGRPCServer.ListAlerts(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
		assert.Nil(t, res)
		mockedAlertService.AssertExpectations(t)
	})
}

func TestGRPCServer_CreateAlertHistory(t *testing.T) {
	timenow := time.Now()

	payload := map[string]interface{}{
		"alerts": []interface{}{
			map[string]interface{}{
				"annotations": map[string]interface{}{
					"metricName":  "bar",
					"metricValue": "30",
					"resource":    "foo",
					"template":    "random",
				},
				"labels": map[string]interface{}{
					"severity": "foo",
				},
				"startsAt": timenow.String(),
				"status":   "foo",
			},
		},
	}

	alertPB := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"status": {
				Kind: &structpb.Value_StringValue{StringValue: "foo"},
			},
			"labels": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"severity": {
								Kind: &structpb.Value_StringValue{StringValue: "foo"},
							},
						},
					},
				},
			},
			"annotations": {
				Kind: &structpb.Value_StructValue{
					StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"resource": {
								Kind: &structpb.Value_StringValue{StringValue: "foo"},
							},
							"template": {
								Kind: &structpb.Value_StringValue{StringValue: "random"},
							},
							"metricName": {
								Kind: &structpb.Value_StringValue{StringValue: "bar"},
							},
							"metricValue": {
								Kind: &structpb.Value_StringValue{StringValue: "30"},
							},
						},
					},
				},
			},
			"startsAt": {
				Kind: &structpb.Value_StringValue{StringValue: timenow.String()},
			},
		},
	}

	dummyReq := &sirenv1beta1.CreateAlertsRequest{
		ProviderId:   1,
		ProviderType: provider.TypeCortex,
		Body: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"alerts": {
					Kind: &structpb.Value_ListValue{
						ListValue: &structpb.ListValue{
							Values: []*structpb.Value{structpb.NewStructValue(alertPB)},
						},
					},
				},
			},
		},
	}

	t.Run("should create alerts objects", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		mockNotificationService := new(mocks.NotificationService)

		dummyAlerts := []alert.Alert{{
			ID:           1,
			ProviderID:   1,
			NamespaceID:  1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Severity:     "CRITICAL",
			Rule:         "random",
			TriggeredAt:  timenow,
		}}
		mockedAlertService.EXPECT().CreateAlerts(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64"), payload).
			Return(dummyAlerts, 1, nil).Once()
		mockNotificationService.EXPECT().Dispatch(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Notification")).Return(nil)

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService, NotificationService: mockNotificationService})

		res, err := dummyGRPCServer.CreateAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertExpectations(t)
	})

	t.Run("should create alerts for resolved alerts", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		mockNotificationService := new(mocks.NotificationService)

		alertPB := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"status": {
					Kind: &structpb.Value_StringValue{StringValue: "resolved"},
				},
				"labels": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"severity": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
							},
						},
					},
				},
				"annotations": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"resource": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
								"template": {
									Kind: &structpb.Value_StringValue{StringValue: "random"},
								},
								"metricName": {
									Kind: &structpb.Value_StringValue{StringValue: "bar"},
								},
								"metricValue": {
									Kind: &structpb.Value_StringValue{StringValue: "30"},
								},
							},
						},
					},
				},
				"startsAt": {
					Kind: &structpb.Value_StringValue{StringValue: timenow.String()},
				},
			},
		}

		dummyReq := &sirenv1beta1.CreateAlertsRequest{
			ProviderId:   1,
			ProviderType: provider.TypeCortex,
			Body: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"alerts": {
						Kind: &structpb.Value_ListValue{
							ListValue: &structpb.ListValue{
								Values: []*structpb.Value{structpb.NewStructValue(alertPB)},
							},
						},
					},
				},
			},
		}

		payload := map[string]interface{}{
			"alerts": []interface{}{
				map[string]interface{}{
					"annotations": map[string]interface{}{
						"metricName":  "bar",
						"metricValue": "30",
						"resource":    "foo",
						"template":    "random",
					},
					"labels": map[string]interface{}{
						"severity": "foo",
					},
					"startsAt": timenow.String(),
					"status":   "resolved",
				},
			},
		}
		dummyAlerts := []alert.Alert{{
			ID:           1,
			ProviderID:   1,
			NamespaceID:  1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Severity:     "resolved",
			Rule:         "random",
			TriggeredAt:  timenow,
		}}
		mockedAlertService.EXPECT().CreateAlerts(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64"), payload).
			Return(dummyAlerts, 1, nil).Once()
		mockNotificationService.EXPECT().Dispatch(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Notification")).Return(nil)

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService, NotificationService: mockNotificationService})

		res, err := dummyGRPCServer.CreateAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "resolved", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertExpectations(t)
	})

	t.Run("should return error Internal if getting alert history failed", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService})

		mockedAlertService.EXPECT().CreateAlerts(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64"), payload).
			Return(nil, 0, errors.New("random error")).Once()

		res, err := dummyGRPCServer.CreateAlerts(context.Background(), dummyReq)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
		assert.Nil(t, res)
		mockedAlertService.AssertExpectations(t)
	})

	t.Run("should insert valid alerts and should not return error if parameters are missing", func(t *testing.T) {
		mockedAlertService := &mocks.AlertService{}
		mockNotificationService := new(mocks.NotificationService)

		payload := map[string]interface{}{
			"alerts": []interface{}{
				map[string]interface{}{
					"annotations": map[string]interface{}{
						"metricName":  "bar",
						"metricValue": "30",
						"resource":    "foo",
						"template":    "random",
					},
					"labels": map[string]interface{}{
						"severity": "foo",
					},
					"startsAt": timenow.String(),
					"status":   "foo",
				}, map[string]interface{}{
					"annotations": map[string]interface{}{
						"metricName":  "bar",
						"metricValue": "30",
						"resource":    "foo",
						"template":    "random",
					},
					"labels": map[string]interface{}{
						"severity": "foo",
					},
					"startsAt": timenow.String(),
					"status":   "resolved",
				},
			},
		}

		alert1PB := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"status": {
					Kind: &structpb.Value_StringValue{StringValue: "foo"},
				},
				"labels": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"severity": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
							},
						},
					},
				},
				"annotations": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"resource": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
								"template": {
									Kind: &structpb.Value_StringValue{StringValue: "random"},
								},
								"metricName": {
									Kind: &structpb.Value_StringValue{StringValue: "bar"},
								},
								"metricValue": {
									Kind: &structpb.Value_StringValue{StringValue: "30"},
								},
							},
						},
					},
				},
				"startsAt": {
					Kind: &structpb.Value_StringValue{StringValue: timenow.String()},
				},
			},
		}

		alert2PB := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"status": {
					Kind: &structpb.Value_StringValue{StringValue: "resolved"},
				},
				"labels": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"severity": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
							},
						},
					},
				},
				"annotations": {
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"resource": {
									Kind: &structpb.Value_StringValue{StringValue: "foo"},
								},
								"template": {
									Kind: &structpb.Value_StringValue{StringValue: "random"},
								},
								"metricName": {
									Kind: &structpb.Value_StringValue{StringValue: "bar"},
								},
								"metricValue": {
									Kind: &structpb.Value_StringValue{StringValue: "30"},
								},
							},
						},
					},
				},
				"startsAt": {
					Kind: &structpb.Value_StringValue{StringValue: timenow.String()},
				},
			},
		}

		dummyReq := &sirenv1beta1.CreateAlertsRequest{
			ProviderId:   1,
			ProviderType: provider.TypeCortex,
			Body: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"alerts": {
						Kind: &structpb.Value_ListValue{
							ListValue: &structpb.ListValue{
								Values: []*structpb.Value{
									structpb.NewStructValue(alert1PB),
									structpb.NewStructValue(alert2PB),
								},
							},
						},
					},
				},
			},
		}

		dummyAlerts := []alert.Alert{{
			ProviderID:   1,
			NamespaceID:  1,
			ResourceName: "foo",
			MetricName:   "bar",
			MetricValue:  "30",
			Rule:         "random",
			Severity:     "CRITICAL",
			TriggeredAt:  time.Now(),
		}}

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), api.HeadersConfig{}, &api.Deps{AlertService: mockedAlertService, NotificationService: mockNotificationService})

		mockedAlertService.EXPECT().CreateAlerts(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("uint64"), mock.AnythingOfType("uint64"), payload).
			Return(dummyAlerts, 2, nil).Once()
		mockNotificationService.EXPECT().Dispatch(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Notification")).Return(nil)

		res, err := dummyGRPCServer.CreateAlerts(context.Background(), dummyReq)
		assert.Equal(t, 1, len(res.GetAlerts()))
		assert.Equal(t, uint64(1), res.GetAlerts()[0].GetProviderId())
		assert.Equal(t, "foo", res.GetAlerts()[0].GetResourceName())
		assert.Equal(t, "random", res.GetAlerts()[0].GetRule())
		assert.Equal(t, "bar", res.GetAlerts()[0].GetMetricName())
		assert.Equal(t, "30", res.GetAlerts()[0].GetMetricValue())
		assert.Equal(t, "CRITICAL", res.GetAlerts()[0].GetSeverity())
		assert.Nil(t, err)
		mockedAlertService.AssertExpectations(t)
	})
}
