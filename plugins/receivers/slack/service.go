package slack

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/plugins/clients/slack"
	"github.com/odpf/siren/store/model"
	"github.com/pkg/errors"
)

type service struct {
	Slacker domain.SlackService
}

func NewService() model.SlackRepository {
	return &service{
		Slacker: nil,
	}
}

var newService = slack.NewService

func (s service) GetWorkspaceChannels(token string) ([]model.Channel, error) {
	s.Slacker = newService(token)
	joinedChannelList, err := s.Slacker.GetJoinedChannelsList()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch joined channel list")
	}

	result := make([]model.Channel, 0)
	for _, c := range joinedChannelList {
		result = append(result, model.Channel{
			ID:   c.ID,
			Name: c.Name,
		})
	}
	return result, nil
}
