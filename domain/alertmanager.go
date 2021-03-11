package domain

type SlackCredential struct {
	Channel string `json:"channel"`
	Webhook string	`json:"webhook"`
	Username string `json:"username"`
}

type SlackConfig struct {
	Critical SlackCredential `json:"critical"`
	Warning SlackCredential  `json:"warning"`
}

type AlertCredential struct {
	Entity               string      `json:"entity"`
	TeamName             string      `json:"team_name"`
	PagerdutyCredentials string      `json:"pagerduty_credentials"`
	SlackConfig          SlackConfig `json:"slack_config"`
}
type AlertmanagerService interface {
	Upsert(credential AlertCredential) (error)
	Get(teamName string)(AlertCredential, error)
	Migrate() error
}
