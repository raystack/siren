package domain

import "github.com/slack-go/slack"

type SlackService interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
	GetUserByEmail(string) (*slack.User, error)
	GetJoinedChannelsList() ([]slack.Channel, error)
}
