package workspace

import "github.com/odpf/siren/domain"

type Channel struct {
	ID   string
	Name string
}

type SlackRepository interface {
	GetWorkspaceChannels(string) ([]Channel, error)
}

func (c *Channel) toDomain() domain.Channel {
	return domain.Channel{
		ID:   c.ID,
		Name: c.Name,
	}
}
