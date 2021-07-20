package slacknotifier

import "github.com/odpf/siren/domain"

type SlackMessage struct {
	ReceiverName string `json:"receiver_name"`
	ReceiverType string `json:"receiver_type"`
	Entity       string `json:"entity"`
	Message      string `json:"message"`
}

func (message *SlackMessage) fromDomain(m *domain.SlackMessage) *SlackMessage {
	message.ReceiverType = m.ReceiverType
	message.ReceiverName = m.ReceiverName
	message.Entity = m.Entity
	message.Message = m.Message
	return message
}

type SlackMessageRepository interface {
	Notify(*SlackMessage, string) error
}
