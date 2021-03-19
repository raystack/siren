package docs

import "github.com/odpf/siren/domain"

// swagger:response templatesResponse
type alertCredentialResponse struct {
	// in:body
	Body domain.AlertCredential
}

//-------------------------
// swagger:route PUT //alertingCredentials/teams/{teamName}:  alertcredential createAlertCredentialRequest
// Upsert AlertCredentials API: This API helps in creating or updating the teams slack and pagerduty credentials
//responses:
//   200: emptyResponse

// swagger:parameters createTemplateRequest
type createAlertCredentialRequest struct {
	// Create AlertCredential request
	// in:body
	Body domain.AlertCredential
}


