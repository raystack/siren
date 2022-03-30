package slack

import (
	"fmt"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/plugins/clients/slack"
	"github.com/pkg/errors"
	goslack "github.com/slack-go/slack"
)

type service struct {
	Slacker domain.SlackService
}

func NewService() *service {
	return &service{
		Slacker: nil,
	}
}

var newService = slack.NewService

func (s *service) GetWorkspaceChannels(token string) ([]Channel, error) {
	s.Slacker = newService(token)
	joinedChannelList, err := s.Slacker.GetJoinedChannelsList()
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
	s.Slacker = newService(token)
	var channelID string
	switch message.ReceiverType {
	case "channel":
		joinedChannelList, err := s.Slacker.GetJoinedChannelsList()
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
		user, err := s.Slacker.GetUserByEmail(message.ReceiverName)
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
	_, _, _, err := s.Slacker.SendMessage(channelID, goslack.MsgOptionText(message.Message, false), goslack.MsgOptionBlocks(message.Blocks.BlockSet...))
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
