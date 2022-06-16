package model

import (
	"time"

	"github.com/odpf/siren/core/alert"
)

type Alert struct {
	ID           uint64 `gorm:"primarykey"`
	Provider     *Provider
	ProviderID   uint64
	ResourceName string
	MetricName   string
	MetricValue  string
	Severity     string
	Rule         string
	TriggeredAt  time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Alert) FromDomain(alrt *alert.Alert) {
	a.ID = alrt.ID
	a.ProviderID = alrt.ProviderID
	a.ResourceName = alrt.ResourceName
	a.MetricName = alrt.MetricName
	a.MetricValue = alrt.MetricValue
	a.Severity = alrt.Severity
	a.Rule = alrt.Rule
	a.TriggeredAt = alrt.TriggeredAt
	a.CreatedAt = alrt.CreatedAt
	a.UpdatedAt = alrt.UpdatedAt
}

func (a *Alert) ToDomain() alert.Alert {
	return alert.Alert{
		ID:           a.ID,
		ProviderID:   a.ProviderID,
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
