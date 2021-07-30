package workspace

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

type ChannelGetter interface {
	GetConversationsForUser(params *slack.GetConversationsForUserParameters) (channels []slack.Channel, nextCursor string, err error)
}

type SlackClient struct {
	Client ChannelGetter
}

var createNewSlackClient = newSlackClient

func newSlackClient(token string) ChannelGetter {
	return slack.New(token)
}

type Repository struct {
	Client ChannelGetter
}

func NewRepository() SlackRepository {
	return &Repository{
		Client: nil,
	}
}

func (r Repository) GetWorkspaceChannel(token string) ([]Channel, error) {
	r.Client = createNewSlackClient(token)
	joinedChannelList, err := getJoinedChannelsList(r.Client)
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

func getJoinedChannelsList(s ChannelGetter) ([]slack.Channel, error) {
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
