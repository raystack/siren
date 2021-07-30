package workspace

import (
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
)

type Repository struct {
	Slacker domain.SlackService
}

func NewRepository(service domain.SlackService) SlackRepository {
	return &Repository{
		Slacker: service,
	}
}

func (r Repository) GetWorkspaceChannel(token string) ([]Channel, error) {
	r.Slacker.UpdateClient(token)
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
