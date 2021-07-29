package slacknotifier

import (
	"errors"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SlackHTTPClientTestSuite struct {
	suite.Suite
}

func TestHTTP(t *testing.T) {
	suite.Run(t, new(SlackHTTPClientTestSuite))
}

func (s *SlackHTTPClientTestSuite) SetupTest() {}

func (s *SlackHTTPClientTestSuite) TestSlackHTTPClient_Notify() {
	s.Run("should notify user identified by their email", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").Return(&slack.User{ID: "U20"}, nil)
		mockedSlackClient.On("SendMessage", "U20",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.Nil(err)
		mockedSlackClient.AssertExpectations(s.T())
	})

	s.Run("should return error if notifying user fails", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").Return(&slack.User{ID: "U20"}, nil)
		mockedSlackClient.On("SendMessage", "U20",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", errors.New("random error"))
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.EqualError(err, "failed to send message to foo@odpf.io: random error")
		mockedSlackClient.AssertExpectations(s.T())
	})

	s.Run("should return error if user lookup by email fails", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetUserByEmail", "foo@odpf.io").Return(nil, errors.New("random error"))
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "foo@odpf.io",
			ReceiverType: "user",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.EqualError(err, "failed to get id for foo@odpf.io: random error")
		mockedSlackClient.AssertExpectations(s.T())
	})

	s.Run("should notify if part of the channel", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("", r.Cursor)
		}).Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "foo",
				Conversation: slack.Conversation{ID: "C01"}},
			}}, "nextCurr", nil).Once()

		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("nextCurr", r.Cursor)
		}).Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "bar",
				Conversation: slack.Conversation{ID: "C02"}},
			}}, "", nil)

		mockedSlackClient.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "foo",
			ReceiverType: "channel",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.Nil(err)
		mockedSlackClient.AssertNumberOfCalls(s.T(), "GetConversationsForUser", 2)
	})

	s.Run("should return error if not part of the channel", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("", r.Cursor)
		}).Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "foo",
				Conversation: slack.Conversation{ID: "C01"}},
			}}, "nextCurr", nil).Once()

		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("nextCurr", r.Cursor)
		}).Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "bar",
				Conversation: slack.Conversation{ID: "C02"}},
			}}, "", nil)

		mockedSlackClient.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "baz",
			ReceiverType: "channel",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.EqualError(err, "app is not part of the channel baz")
	})

	s.Run("should return error failed to fetch joined channels list", func() {
		testNotifierClient := NewSlackNotifierClient()
		oldClientCreator := createNewSlackClient
		defer func() {
			createNewSlackClient = oldClientCreator
		}()
		mockedSlackClient := &SlackClientMock{}
		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("", r.Cursor)
		}).Return([]slack.Channel{
			{GroupConversation: slack.GroupConversation{
				Name:         "foo",
				Conversation: slack.Conversation{ID: "C01"}},
			}}, "nextCurr", nil).Once()

		mockedSlackClient.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("nextCurr", r.Cursor)
		}).Return([]slack.Channel{}, "", errors.New("random error"))

		mockedSlackClient.On("SendMessage", "C01",
			mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
		createNewSlackClient = func(token string) SlackCaller {
			s.Equal("foo_bar", token)
			return mockedSlackClient
		}
		dummyMessage := &SlackMessage{
			ReceiverName: "baz",
			ReceiverType: "channel",
			Message:      "random text",
			Entity:       "odpf",
		}
		err := testNotifierClient.Notify(dummyMessage, "foo_bar")
		s.EqualError(err, "failed to fetch joined channel list: random error")
	})
}
