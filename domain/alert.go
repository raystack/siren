package domain

import "time"

type Alerts struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Id           uint64    `json:"id"`
	ProviderId   uint64    `json:"provider_id"`
	ResourceName string    `json:"resource_name"`
	MetricName   string    `json:"metric_name"`
	MetricValue  string    `json:"metric_value"`
	Severity     string    `json:"severity"`
	Rule         string    `json:"rule"`
	TriggeredAt  time.Time `json:"triggered_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AlertService interface {
	Create(*Alerts) ([]Alert, error)
	Get(string, uint64, uint64, uint64) ([]Alert, error)
	Migrate() error
}
