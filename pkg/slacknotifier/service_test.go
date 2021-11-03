package slacknotifier

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	notifierMock  *MockSlackNotifier
}

func (s *ServiceTestSuite) SetupTest() {
	s.notifierMock = &MockSlackNotifier{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_Notify() {
	dummyMessage := &domain.SlackMessage{
		ReceiverName: "foo",
		ReceiverType: "user",
		Message:      "some text",
		Token:        "token",
	}
	dummyService := Service{
		client:              s.notifierMock,
	}
	s.Run("should call notifier and return success response", func() {
		s.notifierMock.On("Notify", mock.AnythingOfType("*slacknotifier.SlackMessage"), "token").
			Run(func(args mock.Arguments) {
				rarg := args.Get(0)
				s.Require().IsType((*SlackMessage)(nil), rarg)
				r := rarg.(*SlackMessage)
				s.Equal("foo", r.ReceiverName)
				s.Equal("user", r.ReceiverType)
				s.Equal("some text", r.Message)
			}).Return(nil).Once()

		res, err := dummyService.Notify(dummyMessage)
		s.Equal(true, res.OK)
		s.Nil(err)
		s.notifierMock.AssertExpectations(s.T())
	})

	s.Run("should return error response if notifying fails", func() {
		s.notifierMock.On("Notify", mock.AnythingOfType("*slacknotifier.SlackMessage"), "token").
			Return(errors.New("random error")).Once()

		res, err := dummyService.Notify(dummyMessage)
		s.Equal(false, res.OK)
		s.EqualError(err, "could not send notification: random error")
		s.notifierMock.AssertExpectations(s.T())
	})
}

