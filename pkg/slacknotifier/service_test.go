package slacknotifier

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	notifierMock  *MockSlackNotifier
	exchangerMock *mocks.CodeExchangeService
}

func (s *ServiceTestSuite) SetupTest() {
	s.notifierMock = &MockSlackNotifier{}
	s.exchangerMock = &mocks.CodeExchangeService{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_Notify() {
	dummyMessage := &domain.SlackMessage{
		ReceiverName: "foo",
		ReceiverType: "user",
		Message:      "some text",
		Entity:       "odpf",
	}
	dummyService := Service{
		client:              s.notifierMock,
		codeExchangeService: s.exchangerMock,
	}
	s.Run("should call notifier and return success response", func() {
		s.exchangerMock.On("GetToken", "odpf").Return("test_token", nil).Once()
		s.notifierMock.On("Notify", mock.AnythingOfType("*slacknotifier.SlackMessage"), "test_token").
			Run(func(args mock.Arguments) {
				rarg := args.Get(0)
				s.Require().IsType((*SlackMessage)(nil), rarg)
				r := rarg.(*SlackMessage)
				s.Equal("foo", r.ReceiverName)
				s.Equal("user", r.ReceiverType)
				s.Equal("some text", r.Message)
				s.Equal("odpf", r.Entity)
			}).Return(nil).Once()

		res, err := dummyService.Notify(dummyMessage)
		s.Equal(true, res.OK)
		s.Nil(err)
		s.notifierMock.AssertExpectations(s.T())
		s.exchangerMock.AssertExpectations(s.T())
	})

	s.Run("should return error response if notifying fails", func() {
		s.exchangerMock.On("GetToken", "odpf").Return("test_token", nil).Once()
		s.notifierMock.On("Notify", mock.AnythingOfType("*slacknotifier.SlackMessage"), "test_token").
			Return(errors.New("random error")).Once()

		res, err := dummyService.Notify(dummyMessage)
		s.Equal(false, res.OK)
		s.EqualError(err, "could not send notification: random error")
		s.notifierMock.AssertExpectations(s.T())
		s.exchangerMock.AssertExpectations(s.T())
	})

	s.Run("should return error response if getting token fails", func() {
		s.exchangerMock.On("GetToken", "odpf").
			Return("", errors.New("random token")).Once()

		res, err := dummyService.Notify(dummyMessage)
		s.Equal(false, res.OK)
		s.EqualError(err, "could not get token for entity: odpf: random token")
		s.exchangerMock.AssertExpectations(s.T())
	})
}

func (s *ServiceTestSuite) TestService_NewService() {
	s.Run("should return error in service initialization", func() {
		res, err := NewService(nil, "rQvRLU4S6NOtJPDBC0ybemgiU710twcN")
		s.Nil(err)
		s.NotNil(res)
	})

	s.Run("should return error in service initialization", func() {
		res, err := NewService(nil, "abcd")
		s.EqualError(err, `failed to init slack notifier service: failed to create codeexchange repository: random hash should be 32 chars in length`)
		s.Nil(res)
	})
}
