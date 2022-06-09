package model

import (
	"time"

	"github.com/odpf/siren/core/alert"
)

type Alert struct {
	Id           uint64 `gorm:"primarykey"`
	Provider     *Provider
	ProviderId   uint64
	ResourceName string
	MetricName   string
	MetricValue  string
	Severity     string
	Rule         string
	TriggeredAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Alert) FromDomain(alert *alert.Alert) {
	a.Id = alert.Id
	a.ProviderId = alert.ProviderId
	a.ResourceName = alert.ResourceName
	a.MetricName = alert.MetricName
	a.MetricValue = alert.MetricValue
	a.Severity = alert.Severity
	a.Rule = alert.Rule
	a.TriggeredAt = alert.TriggeredAt
	a.CreatedAt = alert.CreatedAt
	a.UpdatedAt = alert.UpdatedAt
}

func (a *Alert) ToDomain() alert.Alert {
	return alert.Alert{
		Id:           a.Id,
		ProviderId:   a.ProviderId,
		ResourceName: a.ResourceName,
		MetricName:   a.MetricName,
		MetricValue:  a.MetricValue,
		Severity:     a.Severity,
		Rule:         a.Rule,
		TriggeredAt:  a.TriggeredAt,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}
