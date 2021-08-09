package workspace

import (
	"errors"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite

	repository    *MockSlackRepository
	exchangerMock mocks.CodeExchangeService
	service       Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.repository = &MockSlackRepository{}
	s.exchangerMock = mocks.CodeExchangeService{}
	s.service = Service{
		client:              s.repository,
		codeExchangeService: &s.exchangerMock,
	}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_GetChannels() {
	s.Run("should return channels on success response", func() {
		s.exchangerMock.On("GetToken", "test_workspace").
			Return("test_token", nil).Once()
		s.repository.On("GetWorkspaceChannels", "test_token").
			Return([]Channel{
				{Name: "foo"},
			}, nil).Once()

		channels, err := s.service.GetChannels("test_workspace")
		s.Equal(1, len(channels))
		s.Equal("foo", channels[0].Name)
		s.Nil(err)
	})

	s.Run("should return error response if get token fail", func() {
		s.exchangerMock.On("GetToken", "test_workspace").
			Return("", errors.New("random error")).Once()

		channels, err := s.service.GetChannels("test_workspace")
		s.Nil(channels)
		s.EqualError(err, "could not get token for workspace: test_workspace: random error")
	})

	s.Run("should return error response if get channels fail", func() {
		s.exchangerMock.On("GetToken", "test_workspace").
			Return("test_token", nil).Once()
		s.repository.On("GetWorkspaceChannels", "test_token").
			Return(nil, errors.New("random error")).Once()

		channels, err := s.service.GetChannels("test_workspace")
		s.Nil(channels)
		s.EqualError(err, "could not get channels for workspace: test_workspace: random error")
	})
}
