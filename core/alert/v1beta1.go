package alert

import (
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *Alert) ToV1beta1Proto() *sirenv1beta1.Alert {
	return &sirenv1beta1.Alert{
		Id:            a.ID,
		ProviderId:    a.ProviderID,
		ResourceName:  a.ResourceName,
		MetricName:    a.MetricName,
		MetricValue:   a.MetricValue,
		Severity:      a.Severity,
		Rule:          a.Rule,
		TriggeredAt:   timestamppb.New(a.TriggeredAt),
		NamespaceId:   a.NamespaceID,
		SilenceStatus: a.SilenceStatus,
		CreatedAt:     timestamppb.New(a.CreatedAt),
		UpdatedAt:     timestamppb.New(a.UpdatedAt),
		GroupKey:      a.GroupKey,
		Status:        a.Status,
		Annotations:   a.Annotations,
		Labels:        a.Labels,
		GeneratorUrl:  a.GeneratorURL,
		Fingerprint:   a.Fingerprint,
	}
}

func FromV1beta1Proto(proto *sirenv1beta1.Alert) *Alert {
	return &Alert{
		ID:            proto.GetId(),
		ProviderID:    proto.GetProviderId(),
		NamespaceID:   proto.GetNamespaceId(),
		ResourceName:  proto.GetResourceName(),
		MetricName:    proto.GetMetricName(),
		MetricValue:   proto.GetMetricValue(),
		Severity:      proto.GetSeverity(),
		Rule:          proto.GetRule(),
		TriggeredAt:   proto.GetTriggeredAt().AsTime(),
		CreatedAt:     proto.GetCreatedAt().AsTime(),
		UpdatedAt:     proto.GetUpdatedAt().AsTime(),
		SilenceStatus: proto.GetSilenceStatus(),

		GroupKey:     proto.GetGroupKey(),
		Status:       proto.GetStatus(),
		Annotations:  proto.GetAnnotations(),
		Labels:       proto.GetLabels(),
		GeneratorURL: proto.GetGeneratorUrl(),
		Fingerprint:  proto.GetFingerprint(),
	}
}
