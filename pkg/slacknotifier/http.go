package slacknotifier

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type SlackCaller interface {
	SendMessage(string, ...slack.MsgOption) (string, string, string, error)
	GetConversations(*slack.GetConversationsParameters) ([]slack.Channel, string, error)
	GetUserByEmail(string) (*slack.User, error)
	JoinConversation(string) (*slack.Channel, string, []string, error)
	GetConversationsForUser(params *slack.GetConversationsForUserParameters) (channels []slack.Channel, nextCursor string, err error)
}

type SlackNotifierClient struct {
	Client SlackCaller
}

func NewSlackNotifierClient() SlackNotifier {
	return &SlackNotifierClient{
		Client: nil,
	}
}

// Notify function takes value receiver because we don't want to share r.client with concurrent requests
func (r SlackNotifierClient) Notify(message *SlackMessage, token string) error {
	r.Client = createNewSlackClient(token)
	return notifyWithClient(message, r.Client)
}

var createNewSlackClient = newSlackClient

func newSlackClient(token string) SlackCaller {
	return slack.New(token)
}

func notifyWithClient(message *SlackMessage, client SlackCaller) error {
	var channelID string
	switch message.ReceiverType {
	case "channel":
		joinedChannelList, err := getJoinedChannelsList(client)
		if err != nil {
			return &JoinedChannelFetchErr{
				Err: errors.Wrap(err, "failed to fetch joined channel list"),
			}
		}
		channelID = searchChannelId(joinedChannelList, message.ReceiverName)
		if channelID == "" {
			return &NoChannelFoundErr{
				Err: errors.New(fmt.Sprintf("app is not part of the channel %s", message.ReceiverName)),
			}
		}
	case "user":
		user, err := client.GetUserByEmail(message.ReceiverName)
		if err != nil {
			if err.Error() == "users_not_found" {
				return &UserLookupByEmailErr{
					Err: errors.Wrap(err, fmt.Sprintf("failed to get id for %s", message.ReceiverName)),
				}
			}
			return &SlackNotifierErr{
				Err: err,
			}
		}
		channelID = user.ID
	}
	_, _, _, err := client.SendMessage(channelID, slack.MsgOptionText(message.Message, false))
	if err != nil {
		return &MsgSendErr{
			Err: errors.Wrap(err, fmt.Sprintf("failed to send message to %s", message.ReceiverName)),
		}
	}
	return nil
}

func getJoinedChannelsList(s SlackCaller) ([]slack.Channel, error) {
	channelList := make([]slack.Channel, 0)
	curr := ""
	for {
		channels, nextCursor, err := s.GetConversationsForUser(&slack.GetConversationsForUserParameters{
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

func searchChannelId(channels []slack.Channel, channelName string) string {
	for _, c := range channels {
		if c.Name == channelName {
			return c.ID
		}
	}
	return ""
}
