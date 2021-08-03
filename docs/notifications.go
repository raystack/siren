package docs

import (
	"github.com/odpf/siren/domain"
	"github.com/slack-go/slack"
)

type SlackMessage struct {
	ReceiverName string `json:"receiver_name" validate:"required"`
	ReceiverType string `json:"receiver_type" validate:"required,receiverTypeChecker"`
	Entity       string `json:"entity" validate:"required"`
	Message      string `json:"message"`
	slack.Blocks
}

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
	Body SlackMessage
}

// POST notifications response
// swagger:response postNotificationsResponse
type postNotificationsResponse struct {
	// in:body
	Body domain.SlackMessageSendResponse
}
