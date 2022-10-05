package v1beta1

import (
	"context"
	"fmt"
	"time"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/provider"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *GRPCServer) ListAlerts(ctx context.Context, req *sirenv1beta1.ListAlertsRequest) (*sirenv1beta1.ListAlertsResponse, error) {
	alerts, err := s.alertService.List(ctx, alert.Filter{
		ResourceName: req.GetResourceName(),
		ProviderID:   req.GetProviderId(),
		StartTime:    int64(req.GetStartTime()),
		EndTime:      int64(req.GetEndTime()),
	})
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

func (s *GRPCServer) CreateCortexAlerts(ctx context.Context, req *sirenv1beta1.CreateCortexAlertsRequest) (*sirenv1beta1.CreateCortexAlertsResponse, error) {
	alerts := make([]*alert.Alert, 0)

	badAlertCount := 0
	firingLen := 0
	// Alert model follows alertmanager webhook contract
	// https://github.com/prometheus/alertmanager/blob/main/notify/webhook/webhook.go#L64
	for _, item := range req.GetAlerts() {

		if item.GetStatus() == "firing" {
			firingLen++
		}

		severity := item.GetLabels()["severity"]
		if item.GetStatus() == "resolved" {
			severity = item.GetStatus()
		}

		alrt := &alert.Alert{
			ProviderID:   req.GetProviderId(),
			ResourceName: fmt.Sprintf("%v", item.GetAnnotations()["resource"]),
			MetricName:   fmt.Sprintf("%v", item.GetAnnotations()["metric_name"]),
			MetricValue:  fmt.Sprintf("%v", item.GetAnnotations()["metric_value"]),
			Severity:     severity,
			Rule:         fmt.Sprintf("%v", item.GetAnnotations()["template"]),
			TriggeredAt:  item.GetStartsAt().AsTime(),
		}
		if !isValidCortexAlert(alrt) {
			badAlertCount++
			continue
		}
		alerts = append(alerts, alrt)
	}
	createdAlerts, err := s.alertService.Create(ctx, alerts)
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

	// Publish to notification service
	for _, a := range req.GetAlerts() {
		variables := map[string]interface{}{}

		for k, v := range a.GetAnnotations() {
			variables[k] = v
		}

		variables["status"] = a.GetStatus()
		variables["generator_url"] = a.GetGeneratorUrl()
		variables["num_alerts_firing"] = firingLen
		// TODO variables["group_key"]

		n := &notification.Notification{
			ProviderType: provider.TypeCortex,
			ID:           "cortex-" + a.GetFingerprint(),
			Variables:    variables,
			Labels:       a.GetLabels(),
			CreatedAt:    time.Now(),
		}

		if err := s.notificationService.Dispatch(ctx, *n); err != nil {
			s.logger.Warn("failed to send to notification service", "api", "alerts", "notification", n, "err", err)
		}
	}

	if badAlertCount > 0 {
		s.logger.Error("parameters are missing for alert", "alert count", badAlertCount)
		return result, nil
	}

	return result, nil
}

func isValidCortexAlert(alrt *alert.Alert) bool {
	return alrt != nil && !(alrt.ResourceName == "" || alrt.Rule == "" ||
		alrt.MetricValue == "" || alrt.MetricName == "" ||
		alrt.Severity == "")
}
