package slackworkspace

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
)

type Repository struct {
	Slacker domain.SlackService
}

func NewRepository() SlackRepository {
	return &Repository{
		Slacker: nil,
	}
}

var newService = slack.NewService

func (r Repository) GetWorkspaceChannels(token string) ([]Channel, error) {
	r.Slacker = newService(token)
	joinedChannelList, err := r.Slacker.GetJoinedChannelsList()
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
