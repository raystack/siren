package workspace

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryTestSuite struct {
	suite.Suite
	repository SlackRepository
	slacker    *mocks.SlackService
}

func TestHTTP(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestGetWorkspaceChannel() {
	oldServiceCreator := newService
	mockedSlackService := &mocks.SlackService{}
	newService = func(string) domain.SlackService {
		return mockedSlackService
	}
	defer func() { newService = oldServiceCreator }()
	s.slacker = mockedSlackService
	s.repository = &Repository{
		Slacker: s.slacker,
	}

	s.Run("should return joined channel list in a workspace", func() {
		s.slacker.On("GetJoinedChannelsList").Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{Name: "foo"}},
			{GroupConversation: slack.GroupConversation{Name: "bar"}}}, nil).Once()
		channels, err := s.repository.GetWorkspaceChannel("test_token")
		s.Equal(2, len(channels))
		s.Equal("foo", channels[0].Name)
		s.Equal("bar", channels[1].Name)
		s.Nil(err)
		s.slacker.AssertExpectations(s.T())
	})

	s.Run("should return error if get joined channel list fail", func() {
		s.slacker.On("GetJoinedChannelsList").
			Return(nil, errors.New("random error")).Once()

		channels, err := s.repository.GetWorkspaceChannel("test_token")
		s.Nil(channels)
		s.EqualError(err, "failed to fetch joined channel list: random error")
	})
}
