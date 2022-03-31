package slack

import (
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	goslack "github.com/slack-go/slack"
)

type slackCaller interface {
	SendMessage(string, ...goslack.MsgOption) (string, string, string, error)
	GetUserByEmail(string) (*goslack.User, error)
	GetConversationsForUser(params *goslack.GetConversationsForUserParameters) (channels []goslack.Channel, nextCursor string, err error)
}

type service struct{}

func NewService() *service {
	return &service{}
}

var newClient = newSlackClient

func (s *service) GetWorkspaceChannels(token string) ([]Channel, error) {
	client := newClient(token)
	joinedChannelList, err := getJoinedChannelsList(client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch joined channel list")
	}

	result := make([]Channel, 0)
	for _, c := range joinedChannelList {
		result = append(result, Channel{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return result, nil
}

func (s *service) Notify(message *domain.SlackMessage) (*domain.SlackMessageSendResponse, error) {
	payload := new(SlackMessage)
	payload.fromDomain(message)

	res := &domain.SlackMessageSendResponse{
		OK: false,
	}

	if err := s.notifyWithClient(payload, message.Token); err != nil {
		return res, errors.Wrap(err, "could not send notification")
	}

	res.OK = true
	return res, nil
}

func (s *service) notifyWithClient(message *SlackMessage, token string) error {
	client := newClient(token)
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
	_, _, _, err := client.SendMessage(channelID, goslack.MsgOptionText(message.Message, false), goslack.MsgOptionBlocks(message.Blocks.BlockSet...))
	if err != nil {
		return &MsgSendErr{
			Err: errors.Wrap(err, fmt.Sprintf("failed to send message to %s", message.ReceiverName)),
		}
	}
	return nil
}

func getJoinedChannelsList(client slackCaller) ([]goslack.Channel, error) {
	channelList := make([]goslack.Channel, 0)
	curr := ""
	for {
		channels, nextCursor, err := client.GetConversationsForUser(&goslack.GetConversationsForUserParameters{
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

func searchChannelId(channels []goslack.Channel, channelName string) string {
	for _, c := range channels {
		if c.Name == channelName {
			return c.ID
		}
	}
	return ""
}

func newSlackClient(token string) slackCaller {
	return goslack.New(token)
}
