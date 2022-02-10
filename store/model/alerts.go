package model

import (
	"github.com/odpf/siren/domain"
	"time"
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

type AlertRepository interface {
	Create(*Alert) (*Alert, error)
	Get(string, uint64, uint64, uint64) ([]Alert, error)
	Migrate() error
}

func (a *Alert) FromDomain(alert *domain.Alert) {
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

func (a *Alert) ToDomain() domain.Alert {
	return domain.Alert{
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
