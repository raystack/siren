package docs

import "github.com/odpf/siren/domain"

type AlertCredentialResponse struct {
	Entity               string             `json:"entity"`
	PagerdutyCredentials string             `json:"pagerduty_credentials"`
	SlackConfig          domain.SlackConfig `json:"slack_config"`
}

//-------------------------
// swagger:route PUT /alertingCredentials/teams/{teamName}  alertcredential createAlertCredentialRequest
// Upsert AlertCredentials API: This API helps in creating or updating the teams slack and pagerduty credentials
//responses:
//   200:

// swagger:parameters createAlertCredentialRequest
type createAlertCredentialRequest struct {
	// Create AlertCredential request
	// in:body
	Body AlertCredentialResponse
	// name of the team
	// in:path
	TeamName string `json:"teamName"`
}

//-------------------------
// swagger:route GET /alertingCredentials/teams/{teamName}  alertcredential getAlertCredentialRequest
// Get AlertCredentials API: This API helps in getting the teams slack and pagerduty credentials
//responses:
//   200: AlertCredentialResponse

// swagger:parameters getAlertCredentialRequest
type getAlertCredentialRequest struct {
	// name of the team
	// in:path
	TeamName string `json:"teamName"`
}
