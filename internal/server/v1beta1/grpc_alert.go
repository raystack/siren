package v1beta1

import (
	"context"

	"github.com/odpf/siren/core/alert"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name=AlertService -r --case underscore --with-expecter --structname AlertService --filename alert_service.go --output=./mocks
type AlertService interface {
	Create(*alert.Alerts) ([]alert.Alert, error)
	Get(string, uint64, uint64, uint64) ([]alert.Alert, error)
}

func (s *GRPCServer) ListAlerts(_ context.Context, req *sirenv1beta1.ListAlertsRequest) (*sirenv1beta1.ListAlertsResponse, error) {
	resourceName := req.GetResourceName()
	providerId := req.GetProviderId()
	startTime := req.GetStartTime()
	endTime := req.GetEndTime()

	alerts, err := s.alertService.Get(resourceName, providerId, startTime, endTime)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Alert{}
	for _, alert := range alerts {
		item := &sirenv1beta1.Alert{
			Id:           alert.ID,
			ProviderId:   alert.ProviderID,
			ResourceName: alert.ResourceName,
			MetricName:   alert.MetricName,
			MetricValue:  alert.MetricValue,
			Severity:     alert.Severity,
			Rule:         alert.Rule,
			TriggeredAt:  timestamppb.New(alert.TriggeredAt),
		}
		items = append(items, item)
	}
	return &sirenv1beta1.ListAlertsResponse{
		Alerts: items,
	}, nil
}

func (s *GRPCServer) CreateCortexAlerts(_ context.Context, req *sirenv1beta1.CreateCortexAlertsRequest) (*sirenv1beta1.CreateCortexAlertsResponse, error) {
	alerts := alert.Alerts{Alerts: make([]alert.Alert, 0)}

	badAlertCount := 0
	for _, item := range req.GetAlerts() {
		severity := item.Labels.GetSeverity()
		if item.GetStatus() == "resolved" {
			severity = item.GetStatus()
		}

		alert := alert.Alert{
			ProviderID:   req.GetProviderId(),
			ResourceName: item.GetAnnotations().GetResource(),
			MetricName:   item.GetAnnotations().GetMetricName(),
			MetricValue:  item.GetAnnotations().GetMetricValue(),
			Severity:     severity,
			Rule:         item.GetAnnotations().GetTemplate(),
			TriggeredAt:  item.GetStartsAt().AsTime(),
		}
		if !isValidCortexAlert(alert) {
			badAlertCount++
			continue
		}
		alerts.Alerts = append(alerts.Alerts, alert)
	}
	createdAlerts, err := s.alertService.Create(&alerts)
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	items := []*sirenv1beta1.Alert{}
	for _, item := range createdAlerts {
		alertHistoryItem := &sirenv1beta1.Alert{
			Id:           item.ID,
			ProviderId:   item.ProviderID,
			ResourceName: item.ResourceName,
			MetricName:   item.MetricName,
			MetricValue:  item.MetricValue,
			Severity:     item.Severity,
			Rule:         item.Rule,
			TriggeredAt:  timestamppb.New(item.TriggeredAt),
		}
		items = append(items, alertHistoryItem)
	}
	result := &sirenv1beta1.CreateCortexAlertsResponse{
		Alerts: items,
	}

	if badAlertCount > 0 {
		s.logger.Error("parameters are missing for alert", "alert count", badAlertCount)
		return result, nil
	}
	return result, nil
}

func isValidCortexAlert(alert alert.Alert) bool {
	return !(alert.ResourceName == "" || alert.Rule == "" ||
		alert.MetricValue == "" || alert.MetricName == "" ||
		alert.Severity == "")
}
