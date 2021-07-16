package codeexchange

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
)

type ServiceTestSuite struct {
	suite.Suite
	exchangerMock  *MockExchanger
	repositoryMock *MockExchangeRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.exchangerMock = &MockExchanger{}
	s.repositoryMock = &MockExchangeRepository{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_Exchange() {
	s.Run("should call repository Upsert method", func() {
		clientID, clientSecret := "test-client-id", "test-client-secret"

		dummyService := Service{
			repository:   s.repositoryMock,
			exchanger:    s.exchangerMock,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		dummyOAuthPayload := domain.OAuthPayload{
			Code:      "test-client-code",
			Workspace: "test-workspace",
		}

		codeExchangeHTTPResponse := CodeExchangeHTTPResponse{
			AccessToken: "test-access-token",
			Team: struct {
				Name string `json:"name"`
			}{
				Name: "test-name",
			},
		}

		s.exchangerMock.On("Exchange", dummyOAuthPayload.Code, clientID, clientSecret).
			Return(codeExchangeHTTPResponse, nil).Once()

		s.repositoryMock.On("Upsert", &AccessToken{
			AccessToken: codeExchangeHTTPResponse.AccessToken,
			Workspace:   dummyOAuthPayload.Workspace,
		}).Return(nil).Once()

		res, err := dummyService.Exchange(dummyOAuthPayload)

		s.Nil(err)
		s.Equal(true, res.OK)
		s.exchangerMock.AssertCalled(s.T(), "Exchange", dummyOAuthPayload.Code, clientID, clientSecret)
		s.repositoryMock.AssertCalled(s.T(), "Upsert", &AccessToken{
			AccessToken: codeExchangeHTTPResponse.AccessToken,
			Workspace:   dummyOAuthPayload.Workspace,
		})
	})

	s.Run("should handle exchange errors", func() {
		clientID, clientSecret := "test-client-id", "test-client-secret"
		exchangerMock := &MockExchanger{}
		repositoryMock := &MockExchangeRepository{}
		dummyService := Service{
			repository:   repositoryMock,
			exchanger:    exchangerMock,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		dummyOAuthPayload := domain.OAuthPayload{
			Code:      "test-client-code",
			Workspace: "test-workspace",
		}

		exchangerMock.On("Exchange", dummyOAuthPayload.Code, clientID, clientSecret).
			Return(CodeExchangeHTTPResponse{}, errors.New("random error")).Once()

		res, err := dummyService.Exchange(dummyOAuthPayload)
		s.Nil(res)
		s.EqualError(err, "failed to exchange code with slack OAuth server: random error")
		exchangerMock.AssertCalled(s.T(), "Exchange", dummyOAuthPayload.Code, clientID, clientSecret)
	})

	s.Run("should handle repository errors", func() {
		clientID, clientSecret := "test-client-id", "test-client-secret"

		dummyService := Service{
			repository:   s.repositoryMock,
			exchanger:    s.exchangerMock,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		dummyOAuthPayload := domain.OAuthPayload{
			Code:      "test-client-code",
			Workspace: "test-workspace",
		}

		codeExchangeHTTPResponse := CodeExchangeHTTPResponse{
			AccessToken: "test-access-token",
			Team: struct {
				Name string `json:"name"`
			}{
				Name: "test-name",
			},
		}

		s.exchangerMock.On("Exchange", dummyOAuthPayload.Code, clientID, clientSecret).
			Return(codeExchangeHTTPResponse, nil).Once()

		s.repositoryMock.On("Upsert", &AccessToken{
			AccessToken: codeExchangeHTTPResponse.AccessToken,
			Workspace:   dummyOAuthPayload.Workspace,
		}).Return(errors.New("random error")).Once()

		res, err := dummyService.Exchange(dummyOAuthPayload)
		s.Nil(res)
		s.EqualError(err, "random error")
	})

	s.Run("should use workspace from exchange response if not provided in user input", func() {
		clientID, clientSecret := "test-client-id", "test-client-secret"

		dummyService := Service{
			repository:   s.repositoryMock,
			exchanger:    s.exchangerMock,
			clientID:     clientID,
			clientSecret: clientSecret,
		}

		dummyOAuthPayload := domain.OAuthPayload{
			Code: "test-client-code",
		}

		codeExchangeHTTPResponse := CodeExchangeHTTPResponse{
			AccessToken: "test-access-token",
			Team: struct {
				Name string `json:"name"`
			}{
				Name: "test-name",
			},
		}

		s.exchangerMock.On("Exchange", dummyOAuthPayload.Code, clientID, clientSecret).
			Return(codeExchangeHTTPResponse, nil).Once()

		s.repositoryMock.On("Upsert", &AccessToken{
			AccessToken: codeExchangeHTTPResponse.AccessToken,
			Workspace:   "test-name",
		}).Return(nil).Once()

		res, err := dummyService.Exchange(dummyOAuthPayload)

		s.Nil(err)
		s.Equal(true, res.OK)
		s.repositoryMock.AssertCalled(s.T(), "Upsert", &AccessToken{
			AccessToken: codeExchangeHTTPResponse.AccessToken,
			Workspace:   "test-name",
		})
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &MockExchangeRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
