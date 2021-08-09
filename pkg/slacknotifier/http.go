package slacknotifier

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
	goslack "github.com/slack-go/slack"
)

type SlackNotifierClient struct {
	Slacker domain.SlackService
}

func NewSlackNotifierClient() SlackNotifier {
	return &SlackNotifierClient{
		Slacker: nil,
	}
}

// Notify function takes value receiver because we don't want to share r.client with concurrent requests
func (r SlackNotifierClient) Notify(message *SlackMessage, token string) error {
	return r.notifyWithClient(message, token)
}

var newService = slack.NewService

func (r SlackNotifierClient) notifyWithClient(message *SlackMessage, token string) error {
	r.Slacker = newService(token)
	var channelID string
	switch message.ReceiverType {
	case "channel":
		joinedChannelList, err := r.Slacker.GetJoinedChannelsList()
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
		user, err := r.Slacker.GetUserByEmail(message.ReceiverName)
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
	_, _, _, err := r.Slacker.SendMessage(channelID, goslack.MsgOptionText(message.Message, false))
	if err != nil {
		return &MsgSendErr{
			Err: errors.Wrap(err, fmt.Sprintf("failed to send message to %s", message.ReceiverName)),
		}
	}
	return nil
}

func searchChannelId(channels []goslack.Channel, channelName string) string {
	for _, c := range channels {
		if c.Name == channelName {
			return c.ID
		}
	}
	return ""
}
