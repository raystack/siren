package v1beta1

import (
	"context"
	"fmt"
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
	items, err := s.createAlerts(ctx, req.GetProviderType(), req.GetProviderId(), 0, req.GetBody().AsMap())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateAlertsResponse{
		Alerts: items,
	}, nil
}

func (s *GRPCServer) CreateAlertsWithNamespace(ctx context.Context, req *sirenv1beta1.CreateAlertsWithNamespaceRequest) (*sirenv1beta1.CreateAlertsWithNamespaceResponse, error) {
	items, err := s.createAlerts(ctx, req.GetProviderType(), req.GetProviderId(), req.GetNamespaceId(), req.GetBody().AsMap())
	if err != nil {
		return nil, s.generateRPCErr(err)
	}

	return &sirenv1beta1.CreateAlertsWithNamespaceResponse{
		Alerts: items,
	}, nil
}

func (s *GRPCServer) createAlerts(ctx context.Context, providerType string, providerID uint64, namespaceID uint64, body map[string]interface{}) ([]*sirenv1beta1.Alert, error) {
	createdAlerts, firingLen, err := s.alertService.CreateAlerts(ctx, providerType, providerID, namespaceID, body)
	if err != nil {
		return nil, err
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

	if len(createdAlerts) > 0 {
		// Publish to notification service
		n := AlertsToNotification(createdAlerts, firingLen, time.Now())

		if err := s.notificationService.DispatchToSubscribers(ctx, namespaceID, n); err != nil {
			s.logger.Warn("failed to send alert as notification", "err", err, "notification", n)
		}
	} else {
		s.logger.Warn("failed to send alert a as notification, empty created alerts")
	}

	return items, nil
}

// Transform alerts and populate Data and Labels to be interpolated to the system-default template
// .Data
// - id
// - status "FIRING"/"RESOLVED"
// - resource
// - template
// - metric_value
// - metric_name
// - generatorUrl
// - numAlertsFiring
// - dashboard
// - playbook
// - summary
// .Labels
// - severity "WARNING"/"CRITICAL"
// - alertname
// - (others labels defined in rules)
func AlertsToNotification(
	as []alert.Alert,
	firingLen int,
	createdTime time.Time,
) notification.Notification {
	sampleAlert := as[0]
	id := "cortex-" + sampleAlert.Fingerprint

	data := map[string]interface{}{}

	mergedAnnotations := map[string][]string{}
	for _, a := range as {
		for k, v := range a.Annotations {
			mergedAnnotations[k] = append(mergedAnnotations[k], v)
		}
	}

	// make unique
	for k, v := range mergedAnnotations {
		mergedAnnotations[k] = removeDuplicateStringValues(v)
	}

	// render annotations
	for k, vSlice := range mergedAnnotations {
		for _, v := range vSlice {
			if _, ok := data[k]; ok {
				data[k] = fmt.Sprintf("%s\n%s", data[k], v)
			} else {
				data[k] = v
			}
		}
	}

	data["status"] = sampleAlert.Status
	data["generatorUrl"] = sampleAlert.GeneratorURL
	data["numAlertsFiring"] = firingLen
	data["id"] = id

	labels := map[string]string{}

	for _, a := range as {
		for k, v := range a.Labels {
			labels[k] = v
		}
	}

	return notification.Notification{
		ID:        id,
		Data:      data,
		Labels:    labels,
		Template:  template.ReservedName_SystemDefault,
		CreatedAt: createdTime,
	}
}

func removeDuplicateStringValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, v := range strSlice {
		if _, value := keys[v]; !value {
			keys[v] = true
			list = append(list, v)
		}
	}
	return list
}
