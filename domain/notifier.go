package domain

import (
	"errors"
	"github.com/slack-go/slack"
	"gopkg.in/go-playground/validator.v9"
)

type NotifierServices struct {
	Slack SlackNotifierService
}

type SlackMessageSendResponse struct {
	OK bool `json:"ok"`
}

type SlackNotifierService interface {
	Notify(*SlackMessage) (*SlackMessageSendResponse, error)
}

type SlackMessage struct {
	ReceiverName string       `json:"receiver_name" validate:"required"`
	ReceiverType string       `json:"receiver_type" validate:"required,oneof=user channel"`
	Entity       string       `json:"entity" validate:"required"`
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
