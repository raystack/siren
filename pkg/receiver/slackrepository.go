package receiver

import (
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/slack"
	"github.com/odpf/siren/store/model"
	"github.com/pkg/errors"
)

type slackRepository struct {
	Slacker domain.SlackService
}

func NewSlackRepository() model.SlackRepository {
	return &slackRepository{
		Slacker: nil,
	}
}

var newService = slack.NewService

func (r slackRepository) GetWorkspaceChannels(token string) ([]model.Channel, error) {
	r.Slacker = newService(token)
	joinedChannelList, err := r.Slacker.GetJoinedChannelsList()
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
