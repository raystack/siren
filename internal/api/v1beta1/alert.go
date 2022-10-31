package v1beta1

import (
	"context"
	"time"

	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/template"
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

func (s *GRPCServer) CreateAlerts(ctx context.Context, req *sirenv1beta1.CreateAlertsRequest) (*sirenv1beta1.CreateAlertsResponse, error) {
	createdAlerts, firingLen, err := s.alertService.CreateAlerts(ctx, req.GetProviderType(), req.GetProviderId(), req.GetBody().AsMap())
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

	result := &sirenv1beta1.CreateAlertsResponse{
		Alerts: items,
	}

	// Publish to notification service
	for _, a := range createdAlerts {
		n := AlertPBToNotification(a, firingLen, a.GroupKey, time.Now())

		if err := s.notificationService.DispatchToSubscribers(ctx, n); err != nil {
			s.logger.Warn("failed to send alert as notification", "err", err, "notification", n)
		}
	}

	return result, nil
}

// Transform alerts and populate Data and Labels to be interpolated to the system-default template
// .Data
// - id
// - status "FIRING"/"RESOLVED"
// - resource
// - template
// - metricValue
// - metricName
// - generatorUrl
// - numAlertsFiring
// - dashboard
// - playbook
// - summary
// .Labels
// - severity "WARNING"/"CRITICAL"
// - alertname
// - (others labels defined in rules)
func AlertPBToNotification(
	a *alert.Alert,
	firingLen int,
	groupKey string,
	createdTime time.Time,
) notification.Notification {
	id := "cortex-" + a.Fingerprint

	data := map[string]interface{}{}

	for k, v := range a.Annotations {
		data[k] = v
	}

	data["status"] = a.Status
	data["generatorUrl"] = a.GeneratorURL
	data["numAlertsFiring"] = firingLen
	data["groupKey"] = groupKey
	data["id"] = id

	return notification.Notification{
		ID:        id,
		Data:      data,
		Labels:    a.Labels,
		Template:  template.ReservedName_SystemDefault,
		CreatedAt: createdTime,
	}
}
