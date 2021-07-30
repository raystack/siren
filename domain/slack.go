package domain

import "github.com/slack-go/slack"

type SlackService interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
	GetConversations(*slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetUserByEmail(string) (*slack.User, error)
	JoinConversation(string) (*slack.Channel, string, []string, error)
	GetConversationsForUser(*slack.GetConversationsForUserParameters) ([]slack.Channel, string, error)
	GetJoinedChannelsList() ([]slack.Channel, error)
	UpdateClient(token string)
}
