package slack

import (
	"github.com/odpf/siren/domain"
	"github.com/slack-go/slack"
)

type SlackMessage struct {
	ReceiverName string       `json:"receiver_name"`
	ReceiverType string       `json:"receiver_type"`
	Message      string       `json:"message"`
	Blocks       slack.Blocks `json:"block"`
}

func (message *SlackMessage) fromDomain(m *domain.SlackMessage) {
	message.ReceiverType = m.ReceiverType
	message.ReceiverName = m.ReceiverName
	message.Message = m.Message
	message.Blocks = m.Blocks
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackRepository interface {
	GetWorkspaceChannels(string) ([]Channel, error)
	Notify(*domain.SlackMessage) (*domain.SlackMessageSendResponse, error)
}
