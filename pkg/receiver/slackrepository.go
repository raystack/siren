package receiver

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/slack"
	"github.com/pkg/errors"
)

type slackRepository struct {
	Slacker domain.SlackService
}

func NewSlackRepository() SlackRepository {
	return &slackRepository{
		Slacker: nil,
	}
}

var newService = slack.NewService

func (r slackRepository) GetWorkspaceChannels(token string) ([]Channel, error) {
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
