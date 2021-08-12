package domain

type SlackCredential struct {
	Channel  string `json:"channel" validate:"required"`
}

type SlackConfig struct {
	Critical SlackCredential `json:"critical" validate:"required,dive,required"`
	Warning  SlackCredential `json:"warning" validate:"required,dive,required"`
}

type AlertCredential struct {
	Entity               string      `json:"entity" validate:"required"`
	TeamName             string      `json:"team_name"`
	PagerdutyCredentials string      `json:"pagerduty_credentials" validate:"required"`
	SlackConfig          SlackConfig `json:"slack_config" validate:"required,dive,required"`
}
type AlertmanagerService interface {
	Upsert(credential AlertCredential) error
	Get(teamName string) (AlertCredential, error)
	Migrate() error
}
