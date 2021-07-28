package slacknotifier

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type SlackClient interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
	GetConversations(*slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetUserByEmail(string) (*slack.User, error)
	JoinConversation(string) (*slack.Channel, string, []string, error)
	GetConversationsForUser(params *slack.GetConversationsForUserParameters) (channels []slack.Channel, nextCursor string, err error)
}

type SlackHTTPClient struct {
	client SlackClient
}

func NewSlackHTTPClient() SlackNotifier {
	return &SlackHTTPClient{
		client: nil,
	}
}

func getJoinedChannelsList(s SlackClient) ([]slack.Channel, error) {
	channelList := make([]slack.Channel, 0)
	curr := ""
	for {
		channels, nextCursor, err := s.GetConversationsForUser(&slack.GetConversationsForUserParameters{
			Types:  []string{"public_channel", "private_channel"},
			Cursor: curr,
			Limit:  1000})
		if err != nil {
			return channelList, errors.Wrap(err, "failed to get joined channels list")
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

// Notify function takes value receiver because we don't want to share r.client with concurrent requests
func (r SlackHTTPClient) Notify(message *SlackMessage, token string) error {
	r.client = slack.New(token)
	var channelID string
	switch message.ReceiverType {
	case "channel":
		joinedChannelList, err := getJoinedChannelsList(r.client)
		if err != nil {
			return errors.Wrap(err, "failed to fetch joined channel list")
		}
		channelID = searchChannelId(joinedChannelList, message.ReceiverName)
		if channelID == "" {
			return errors.New(fmt.Sprintf("app is not part of the channel %s", message.ReceiverName))
		}
	case "user":
		user, err := r.client.GetUserByEmail(message.ReceiverName)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to get id for %s", message.ReceiverName))
		}
		channelID = user.ID
	}
	_, _, _, err := r.client.SendMessage(channelID, slack.MsgOptionText(message.Message, false))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to send message to %s", message.ReceiverName))
	}
	return nil
}

func searchChannelId(channels []slack.Channel, channelName string) string {
	for _, c := range channels {
		if c.Name == channelName {
			return c.ID
		}
	}
	return ""
}
