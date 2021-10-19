package receiver

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ServiceTestSuite struct {
	suite.Suite
	slackHelperMock *MockSlackHelper
	repositoryMock  *MockReceiverRepository
	slacker         *MockSlackRepository
}

func (s *ServiceTestSuite) SetupTest() {
	s.slackHelperMock = &MockSlackHelper{}
	s.repositoryMock = &MockReceiverRepository{}
	s.slacker = &MockSlackRepository{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_ListReceivers() {
	configurations := make(StringInterfaceMap)
	labels := make(StringStringMap)
	labels["foo"] = "bar"

	s.Run("should call repository List method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock}
		dummyReceivers := []*domain.Receiver{
			{
				Id:             10,
				Name:           "foo",
				Type:           "slack",
				Labels:         labels,
				Configurations: configurations,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		receivers := []*Receiver{
			{
				Id:             10,
				Name:           "foo",
				Type:           "slack",
				Labels:         labels,
				Configurations: configurations,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		s.repositoryMock.On("List").Return(receivers, nil).Once()

		result, err := dummyService.ListReceivers()
		s.Nil(err)
		s.Equal(len(dummyReceivers), len(result))
		s.Equal(dummyReceivers[0].Name, result[0].Name)
		s.repositoryMock.AssertCalled(s.T(), "List")
	})

	s.Run("should call repository List method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock}
		s.repositoryMock.On("List").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.ListReceivers()
		s.Nil(result)
		s.EqualError(err, "service.repository.List: random error")
		s.repositoryMock.AssertCalled(s.T(), "List")
	})
}

func (s *ServiceTestSuite) TestService_CreateReceiver() {
	configurations := make(StringInterfaceMap)
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()

	receiverRequest := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	receiver := &Receiver{
		Id:     10,
		Name:   "foo",
		Type:   "slack",
		Labels: labels,
		Configurations: map[string]interface{}{
			"workspace": "test-name",
			"token":     "token",
		},
		CreatedAt: timenow,
		UpdatedAt: timenow,
	}

	receiverResponse := &domain.Receiver{
		Id:     10,
		Name:   "foo",
		Type:   "slack",
		Labels: labels,
		Configurations: map[string]interface{}{
			"workspace": "test-name",
			"token":     "token",
		},
		CreatedAt: timenow,
		UpdatedAt: timenow,
	}

	s.Run("should call repository Create method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(receiverResponse, nil).Once()
		s.repositoryMock.On("Create", receiver).Return(receiver, nil).Once()
		result, err := dummyService.CreateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(receiverResponse, result)
		s.repositoryMock.AssertCalled(s.T(), "Create", receiver)
	})

	s.Run("should call repository Create method and return error if transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.CreateReceiver(receiverRequest)
		s.Nil(result)
		s.EqualError(err, "slackHelper.Transform: random error")
	})

	s.Run("should call repository Create method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(receiverResponse, nil).Once()
		s.repositoryMock.On("Create", receiver).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.CreateReceiver(receiverRequest)
		s.Nil(result)
		s.EqualError(err, "service.repository.Create: random error")
		s.repositoryMock.AssertCalled(s.T(), "Create", receiver)
	})
}

func (s *ServiceTestSuite) TestService_GetReceiver() {
	receiverID := uint64(10)
	configurations := make(StringInterfaceMap)
	configurations["token"] = "key"

	labels := make(StringStringMap)
	labels["foo"] = "bar"

	data := make(StringInterfaceMap)
	data["channels"] = "[{\"id\":\"1\",\"name\":\"foo\"}]"

	timenow := time.Now()
	dummyReceiver := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		Data:           data,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	receiver := &Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	s.Run("should call repository Get method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackRepository: s.slacker}
		configurations["token"] = "key"
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("Decrypt", "key").
			Return("token", nil).Once()
		s.slacker.On("GetWorkspaceChannels", "token").
			Return([]Channel{
				{ID: "1", Name: "foo"},
			}, nil).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(err)
		s.Equal(dummyReceiver, result)
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should call repository Get method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock}
		newConfigurations := make(StringInterfaceMap)
		newConfigurations["token"] = "key"
		receiver.Configurations = newConfigurations

		s.repositoryMock.On("Get", receiverID).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should return error if slack token decryption failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackRepository: s.slacker}
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("Decrypt", "key").
			Return("", errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "slackHelper.Decrypt: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
		s.slackHelperMock.AssertCalled(s.T(), "Decrypt", "key")
	})

	s.Run("should return error if getting slack channels failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackRepository: s.slacker}
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("Decrypt", "key").
			Return("token", nil).Once()
		s.slacker.On("GetWorkspaceChannels", "token").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "could not get channels: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should return error if invalid slack channels", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackRepository: s.slacker}
		oldjsonMarshal := jsonMarshal
		jsonMarshal = func(interface{}) ([]byte, error) {
			return nil, errors.New("random error")
		}
		defer func() { jsonMarshal = oldjsonMarshal }()

		newConfigurations := make(StringInterfaceMap)
		newConfigurations["token"] = "key"
		receiver.Configurations = newConfigurations

		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("Decrypt", "key").
			Return("token", nil).Once()
		s.slacker.On("GetWorkspaceChannels", "token").
			Return([]Channel{
				{ID: "1", Name: string([]byte{0xff})},
			}, nil).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "invalid channels: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})
}

func (s *ServiceTestSuite) TestService_UpdateReceiver() {
	timenow := time.Now()
	configurations := make(StringInterfaceMap)
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(StringStringMap)
	labels["foo"] = "bar"
	receiverRequest := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	receiver := &Receiver{
		Id:     10,
		Name:   "foo",
		Type:   "slack",
		Labels: labels,
		Configurations: map[string]interface{}{
			"workspace": "test-name",
			"token":     "token",
		},
		CreatedAt: timenow,
		UpdatedAt: timenow,
	}

	receiverResponse := &domain.Receiver{
		Id:     10,
		Name:   "foo",
		Type:   "slack",
		Labels: labels,
		Configurations: map[string]interface{}{
			"workspace": "test-name",
			"token":     "token",
		},
		CreatedAt: timenow,
		UpdatedAt: timenow,
	}

	s.Run("should call repository Update method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(receiverResponse, nil).Once()
		s.repositoryMock.On("Update", receiver).Return(receiver, nil).Once()

		result, err := dummyService.UpdateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(receiverResponse, result)
		s.repositoryMock.AssertCalled(s.T(), "Update", receiver)
	})

	s.Run("should call repository Create method and return error if transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.UpdateReceiver(receiverRequest)
		s.Nil(result)
		s.EqualError(err, "slackHelper.Transform: random error")
	})

	s.Run("should call repository Update method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("Transform", receiverRequest).
			Return(receiverResponse, nil).Once()
		s.repositoryMock.On("Update", receiver).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.UpdateReceiver(receiverRequest)
		s.Nil(result)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Update", receiver)
	})
}

func (s *ServiceTestSuite) TestService_DeleteReceiver() {
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"

	labels := make(StringStringMap)
	labels["foo"] = "bar"
	receiverID := uint64(10)

	s.Run("should call repository Delete method and return nil if no error", func() {
		dummyService := Service{repository: s.repositoryMock}
		s.repositoryMock.On("Delete", receiverID).Return(nil).Once()

		err := dummyService.DeleteReceiver(receiverID)
		s.Nil(err)
		s.repositoryMock.AssertCalled(s.T(), "Delete", receiverID)
	})

	s.Run("should call repository Delete method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock}
		s.repositoryMock.On("Delete", receiverID).
			Return(errors.New("random error")).Once()

		err := dummyService.DeleteReceiver(receiverID)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Delete", receiverID)
	})
}

func (s *ServiceTestSuite) TestService_Migrate() {
	s.Run("should call repository Migrate method and return result", func() {
		dummyService := Service{repository: s.repositoryMock}
		s.repositoryMock.On("Migrate").Return(nil).Once()

		err := dummyService.Migrate()
		s.Nil(err)
		s.repositoryMock.AssertCalled(s.T(), "Migrate")
	})
}
