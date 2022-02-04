package slack

import (
	"github.com/odpf/siren/domain"
	"github.com/slack-go/slack"
)

type SlackCaller interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
	GetUserByEmail(string) (*slack.User, error)
	GetConversationsForUser(params *slack.GetConversationsForUserParameters) (channels []slack.Channel, nextCursor string, err error)
}

type ClientService struct {
	SlackClient SlackCaller
}

func (c *ClientService) GetJoinedChannelsList() ([]slack.Channel, error) {
	channelList := make([]slack.Channel, 0)
	curr := ""
	for {
		channels, nextCursor, err := c.GetConversationsForUser(&slack.GetConversationsForUserParameters{
			Types:  []string{"public_channel", "private_channel"},
			Cursor: curr,
			Limit:  1000})
		if err != nil {
			return channelList, err
		}
		channelList = append(channelList, channels...)
		curr = nextCursor
		if curr == "" {
			break
		}
	}
	return channelList, nil
}

func (c *ClientService) SendMessage(channel string, option ...slack.MsgOption) (string, string, string, error) {
	return c.SlackClient.SendMessage(channel, option...)
}

func (c *ClientService) GetUserByEmail(email string) (*slack.User, error) {
	return c.SlackClient.GetUserByEmail(email)
}

func (c *ClientService) GetConversationsForUser(parameters *slack.GetConversationsForUserParameters) ([]slack.Channel, string, error) {
	return c.SlackClient.GetConversationsForUser(parameters)
}

func NewService(token string) domain.SlackService {
	return &ClientService{SlackClient: slack.New(token)}
}
