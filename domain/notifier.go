package domain

type NotifierServices struct {
	Slack SlackNotifierService
}

type SlackMessage struct {
	ReceiverName string `json:"receiver_name" validate:"required"`
	ReceiverType string `json:"receiver_type" validate:"required,receiverTypeChecker"`
	Entity       string `json:"entity" validate:"required"`
	Message      string `json:"message" validate:"required"`
}

type SlackMessageSendResponse struct {
	OK bool `json:"ok"`
}

type SlackNotifierService interface {
	Notify(*SlackMessage) (*SlackMessageSendResponse, error)
}
