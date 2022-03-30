package slack

import (
	"errors"
	"testing"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SlackRepositoryTestSuite struct {
	suite.Suite
	service SlackRepository
	slacker *mocks.SlackService
}

func TestSlackRepository(t *testing.T) {
	suite.Run(t, new(SlackRepositoryTestSuite))
}

func (s *SlackRepositoryTestSuite) TestGetWorkspaceChannel() {
	oldServiceCreator := newService
	mockedSlackService := &mocks.SlackService{}
	newService = func(string) domain.SlackService {
		return mockedSlackService
	}
	defer func() { newService = oldServiceCreator }()
	s.slacker = mockedSlackService
	s.service = &service{
		Slacker: s.slacker,
	}

	s.Run("should return joined channel list in a workspace", func() {
		s.slacker.On("GetJoinedChannelsList").Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{Name: "foo"}},
			{GroupConversation: slack.GroupConversation{Name: "bar"}}}, nil).Once()
		channels, err := s.service.GetWorkspaceChannels("test_token")
		s.Equal(2, len(channels))
		s.Equal("foo", channels[0].Name)
		s.Equal("bar", channels[1].Name)
		s.Nil(err)
		s.slacker.AssertExpectations(s.T())
	})

	s.Run("should return error if get joined channel list fail", func() {
		s.slacker.On("GetJoinedChannelsList").
			Return(nil, errors.New("random error")).Once()

		channels, err := s.service.GetWorkspaceChannels("test_token")
		s.Nil(channels)
		s.EqualError(err, "failed to fetch joined channel list: random error")
	})
}

func (s *SlackRepositoryTestSuite) TestNotify() {
	oldServiceCreator := newService
	mockedSlackService := &mocks.SlackService{}
	newService = func(string) domain.SlackService {
		return mockedSlackService
	}
	defer func() { newService = oldServiceCreator }()
	s.slacker = mockedSlackService
	s.service = &service{
		Slacker: s.slacker,
	}

	s.Run("should notify user identified by their email", func() {
		mockedSlackService.On("GetUserByEmail", "foo@odpf.io").
			Return(&slack.User{ID: "U20"}, nil).Once()
		mockedSlackService.On("SendMessage", "U20",
			mock.AnythingOfType("slack.MsgOption"), mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()
		dummyMessage := &domain.SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.Nil(err)
		s.True(res.OK)
	})

	s.Run("should return error if notifying user fails", func() {
		mockedSlackService.On("GetUserByEmail", "foo@odpf.io").
			Return(&slack.User{ID: "U20"}, nil).Once()
		mockedSlackService.On("SendMessage", "U20",
			mock.AnythingOfType("slack.MsgOption"),
			mock.AnythingOfType("slack.MsgOption"),
		).Return("", "", "", errors.New("random error")).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.EqualError(err, "could not send notification: failed to send message to foo@odpf.io: random error")
		s.False(res.OK)
	})

	s.Run("should return error if user lookup by email fails", func() {
		mockedSlackService.On("GetUserByEmail", "foo@odpf.io").
			Return(nil, errors.New("users_not_found")).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.EqualError(err, "could not send notification: failed to get id for foo@odpf.io: users_not_found")
		s.False(res.OK)
	})

	s.Run("should return error if user lookup by email returns any error", func() {
		mockedSlackService.On("GetUserByEmail", "foo@odpf.io").
			Return(nil, errors.New("random error")).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.EqualError(err, "could not send notification: random error")
		s.False(res.OK)
	})

	s.Run("should notify if part of the channel", func() {
		mockedSlackService.On("GetJoinedChannelsList").Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "foo",
				Conversation: slack.Conversation{ID: "C01"}},
			}, {GroupConversation: slack.GroupConversation{
				Name:         "bar",
				Conversation: slack.Conversation{ID: "C02"}},
			}}, nil).Once()

		mockedSlackService.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption"),
			mock.AnythingOfType("slack.MsgOption"),
		).Return("", "", "", nil).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "foo",
			ReceiverType: "channel",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.Nil(err)
		s.True(res.OK)
		mockedSlackService.AssertExpectations(s.T())
	})

	s.Run("should return error if not part of the channel", func() {
		mockedSlackService.On("GetJoinedChannelsList").Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "foo",
				Conversation: slack.Conversation{ID: "C01"}},
			}, {GroupConversation: slack.GroupConversation{
				Name:         "bar",
				Conversation: slack.Conversation{ID: "C02"}},
			}}, nil).Once()

		mockedSlackService.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "baz",
			ReceiverType: "channel",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.EqualError(err, "could not send notification: app is not part of the channel baz")
		s.False(res.OK)
	})

	s.Run("should return error failed to fetch joined channels list", func() {
		mockedSlackService.On("GetJoinedChannelsList").
			Return(nil, errors.New("random error")).Once()
		mockedSlackService.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil).Once()

		dummyMessage := &domain.SlackMessage{
			ReceiverName: "baz",
			ReceiverType: "channel",
			Message:      "random text",
			Token:        "foo_bar",
		}
		res, err := s.service.Notify(dummyMessage)
		s.EqualError(err, "could not send notification: failed to fetch joined channel list: random error")
		s.False(res.OK)
	})
}
