package domain

import (
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
	ReceiverName string `json:"receiver_name" validate:"required"`
	ReceiverType string `json:"receiver_type" validate:"required,oneof=user channel"`
	Entity       string `json:"entity" validate:"required"`
	Message      string `json:"message" validate:"required"`
}

func (sm *SlackMessage) Validate() error {
	v := validator.New()
	return v.Struct(sm)
}
