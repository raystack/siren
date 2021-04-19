package alert_history

import (
	"github.com/odpf/siren/domain"
	"strings"
	"time"
)

type Alert struct {
	ID          uint64 `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Resource    string
	Template    string
	MetricName  string
	MetricValue string
	Level       string
}

type AlertHistoryRepository interface {
	Create(*Alert) (*Alert, error)
	Get(string, uint32, uint32) ([]Alert, error)
	Migrate() error
}

func (a *Alert) fromDomain(alert *domain.Alert) {
	a.Template = alert.Annotations.Template
	a.Resource = alert.Annotations.Resource
	a.MetricName = alert.Annotations.MetricName
	a.MetricValue = alert.Annotations.MetricValue
	if alert.Status == "resolved" {
		a.Level = strings.ToUpper(alert.Status)
	} else {
		a.Level = strings.ToUpper(alert.Labels.Severity)
	}
}

func (a *Alert) toDomain() domain.AlertHistoryObject {
	return domain.AlertHistoryObject{
		ID:          a.ID,
		Name:        a.Resource,
		TemplateID:  a.Template,
		MetricValue: a.MetricValue,
		MetricName:  a.MetricName,
		Level:       a.Level,
		Created:     a.CreatedAt,
		Updated:     a.UpdatedAt,
	}
}
