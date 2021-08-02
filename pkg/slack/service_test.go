package slack

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SlackServiceTestSuite struct {
	suite.Suite
	service           ClientService
	mockedSlackCaller *MockSlackCaller
}

func TestHTTP(t *testing.T) {
	suite.Run(t, new(SlackServiceTestSuite))
}

func (s *SlackServiceTestSuite) SetupTest() {
	mockedSlackCaller := &MockSlackCaller{}
	s.mockedSlackCaller = mockedSlackCaller
	s.service = ClientService{
		SlackClient: mockedSlackCaller,
	}
}

func (s *SlackServiceTestSuite) TestService_SendMessage() {
	s.Run("should call SendMessage method", func() {
		s.mockedSlackCaller.On("SendMessage", "test", mock.AnythingOfType("slack.MsgOption")).
			Return("foo", "bar", "baz", errors.New("random error")).Once()
		channel, ts, response, err := s.service.SendMessage("test", slack.MsgOptionText("some text", false))
		s.Equal("foo", channel)
		s.Equal("bar", ts)
		s.Equal("baz", response)
		s.EqualError(err, "random error")
		s.mockedSlackCaller.AssertExpectations(s.T())
	})

}

func (s *SlackServiceTestSuite) TestService_GetUserByEmail() {
	s.Run("should call GetUserByEmail method", func() {
		s.mockedSlackCaller.On("GetUserByEmail", "foo@odpf.io").
			Return(&slack.User{Name: "foo"}, errors.New("random error")).Once()
		user, err := s.service.GetUserByEmail("foo@odpf.io")
		s.Equal("foo", user.Name)
		s.EqualError(err, "random error")
		s.mockedSlackCaller.AssertExpectations(s.T())
	})

}

func (s *SlackServiceTestSuite) TestService_GetConversationsForUser() {
	s.Run("should call GetConversationsForUser method", func() {
		s.mockedSlackCaller.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).
			Return([]slack.Channel{{GroupConversation: slack.GroupConversation{Name: "foo"}},
			}, "test", errors.New("random error")).Once()
		channels, res, err := s.service.GetConversationsForUser(&slack.GetConversationsForUserParameters{})
		s.Equal(1, len(channels))
		s.Equal("foo", channels[0].Name)
		s.Equal("test", res)
		s.EqualError(err, "random error")
		s.mockedSlackCaller.AssertExpectations(s.T())
	})

}

func (s *SlackServiceTestSuite) TestService_GetJoinedChannelsList() {
	s.Run("should return joined channels list", func() {
		s.mockedSlackCaller.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
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

		s.mockedSlackCaller.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
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
			}}, "", nil).Once()

		channels, err := s.service.GetJoinedChannelsList()
		s.Equal(2, len(channels))
		s.Equal("foo", channels[0].Name)
		s.Equal("bar", channels[1].Name)
		s.Nil(err)
		s.mockedSlackCaller.AssertExpectations(s.T())
	})

	s.Run("should return error in getting joined channels list", func() {
		s.mockedSlackCaller.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
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

		s.mockedSlackCaller.On("GetConversationsForUser", mock.AnythingOfType("*slack.GetConversationsForUserParameters")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*slack.GetConversationsForUserParameters)(nil), rarg)
			r := rarg.(*slack.GetConversationsForUserParameters)
			s.Equal(1000, r.Limit)
			s.Equal([]string{"public_channel", "private_channel"}, r.Types)
			s.Equal("nextCurr", r.Cursor)
		}).Return(nil, "", errors.New("random error")).Once()

		channels, err := s.service.GetJoinedChannelsList()
		s.Equal(1, len(channels))
		s.Equal("foo", channels[0].Name)
		s.EqualError(err, "random error")
		s.mockedSlackCaller.AssertExpectations(s.T())
	})
}

func (s *SlackServiceTestSuite) Test_NewService() {
	s.Run("should initialize client with nil", func() {
		service := NewService()
		res, ok := service.(*ClientService)
		s.Equal(true, ok)
		s.Nil(res.SlackClient)
	})
}
