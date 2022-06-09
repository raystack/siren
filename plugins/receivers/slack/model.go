package slack

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"gopkg.in/go-playground/validator.v9"
)

// type SlackMessage struct {
// 	ReceiverName string       `json:"receiver_name"`
// 	ReceiverType string       `json:"receiver_type"`
// 	Message      string       `json:"message"`
// 	Blocks       slack.Blocks `json:"block"`
// }

// func (message *SlackMessage) fromDomain(m *receiver.SlackMessage) {
// 	message.ReceiverType = m.ReceiverType
// 	message.ReceiverName = m.ReceiverName
// 	message.Message = m.Message
// 	message.Blocks = m.Blocks
// }

type SlackMessageSendResponse struct {
	OK bool `json:"ok"`
}

type SlackMessage struct {
	ReceiverName string       `json:"receiver_name" validate:"required"`
	ReceiverType string       `json:"receiver_type" validate:"required,oneof=user channel"`
	Token        string       `json:"token" validate:"required"`
	Message      string       `json:"message"`
	Blocks       slack.Blocks `json:"blocks"`
}

func (sm *SlackMessage) Validate() error {
	v := validator.New()
	if sm.Message == "" && len(sm.Blocks.BlockSet) == 0 {
		return errors.New("non empty message or non zero length block is required")
	}
	return v.Struct(sm)
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SlackRepository interface {
	GetWorkspaceChannels(string) ([]Channel, error)
	Notify(*SlackMessage) (*SlackMessageSendResponse, error)
}
