package domain

type NotifierServices struct {
	Slack SlackNotifierService
}

type SlackMessage struct {
	ReceiverName string `json:"receiver_name"`
	ReceiverType string `json:"receiver_type"`
	Entity       string `json:"entity"`
	Message      string `json:"message"`
}

type SlackMessageSendResponse struct {
	OK bool `json:"ok"`
}

type SlackNotifierService interface {
	Notify(*SlackMessage) (*SlackMessageSendResponse, error)
}
