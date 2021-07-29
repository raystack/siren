package docs

import "github.com/odpf/siren/domain"

//-------------------------
//-------------------------
// swagger:route POST /notifications notifications postNotificationsRequest
// POST Notifications API
// This API sends notifications to configured channel
// responses:
//   200: postResponse

// swagger:parameters postNotificationsRequest
type postNotificationsRequest struct {
	// in:query
	Provider string `json:"provider"`
	// in:body
	Body domain.SlackMessage
}

// POST notifications response
// swagger:response postNotificationsResponse
type postNotificationsResponse struct {
	// in:body
	Body domain.SlackMessageSendResponse
}
