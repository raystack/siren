package slacknotifier

import (
	"fmt"
	"github.com/slack-go/slack"
)

type SlackClient interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
}

type Repository struct {
	slackClient SlackClient
}

func NewRepository() SlackMessageRepository {
	return &Repository{
		slackClient: nil,
	}
}

// Notify function takes value receiver because we don't want to share r.slackClient with concurrent requests
func (r Repository) Notify(message *SlackMessage, token string) error {
	r.slackClient = slack.New(token)
	channel, timestamp, text, err := r.slackClient.SendMessage(message.ReceiverName,
		slack.MsgOptionText(message.Message, false))
	fmt.Println(channel, timestamp, text, err)
	return nil
}

func f(entity string) string {
	return entity
}
