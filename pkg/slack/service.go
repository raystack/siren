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

func (c *ClientService) UpdateClient(token string) {
	c.SlackClient = slack.New(token)
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
		for _, c := range channels {
			channelList = append(channelList, c)
		}
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

func NewService() domain.SlackService {
	return &ClientService{SlackClient: nil}
}
