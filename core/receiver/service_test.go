package receiver

import (
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/plugins/receivers/slack"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	slackHelperMock *MockSlackHelper
	repositoryMock  *MockReceiverRepository
	slacker         *MockSlackService
}

func (s *ServiceTestSuite) SetupTest() {
	s.slackHelperMock = &MockSlackHelper{}
	s.repositoryMock = &MockReceiverRepository{}
	s.slacker = &MockSlackService{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_ListReceivers() {
	configurations := make(map[string]interface{})
	labels := make(map[string]string)
	labels["foo"] = "bar"

	s.Run("should call repository List method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
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
		receivers := []*domain.Receiver{
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
		configurations["token"] = "key"

		s.repositoryMock.On("List").Return(receivers, nil).Once()
		s.slackHelperMock.On("PostTransform", receivers[0]).
			Return(nil).Once()

		result, err := dummyService.ListReceivers()
		s.Nil(err)
		s.Equal(len(dummyReceivers), len(result))
		s.Equal(dummyReceivers[0].Name, result[0].Name)
		s.repositoryMock.AssertCalled(s.T(), "List")
	})

	s.Run("should call repository List method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock}
		configurations["token"] = "key"
		s.repositoryMock.On("List").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.ListReceivers()
		s.Nil(result)
		s.EqualError(err, "service.repository.List: random error")
		s.repositoryMock.AssertCalled(s.T(), "List")
	})

	s.Run("should call repository List method and return error if post slack transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		receivers := []*domain.Receiver{
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
		configurations["token"] = "key"

		s.repositoryMock.On("List").
			Return(receivers, nil).Once()
		s.slackHelperMock.On("PostTransform", receivers[0]).
			Return(errors.New("random error")).Once()

		result, err := dummyService.ListReceivers()
		s.EqualError(err, "slackHelper.PostTransform: random error")
		s.Nil(result)
		s.repositoryMock.AssertCalled(s.T(), "List")
	})
}

func (s *ServiceTestSuite) TestService_CreateReceiver() {
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(map[string]string)
	labels["foo"] = "bar"
	timenow := time.Now()

	tokenConfigurations := map[string]interface{}{
		"workspace": "test-name",
		"token":     "key",
	}
	receiverRequest := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	transfromResponse := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: tokenConfigurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	receiver := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: tokenConfigurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	s.Run("should call repository Create method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*domain.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", receiver).Return(nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(nil).Once()

		err := dummyService.CreateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(transfromResponse, receiverRequest)
		s.repositoryMock.AssertCalled(s.T(), "Create", receiver)
	})

	s.Run("should call repository Create method and return error if pre transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Return(errors.New("random error")).Once()

		err := dummyService.CreateReceiver(receiverRequest)
		s.EqualError(err, "slackHelper.PreTransform: random error")
	})

	s.Run("should call repository Create method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*domain.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", receiver).
			Return(errors.New("random error")).Once()

		err := dummyService.CreateReceiver(receiverRequest)
		s.EqualError(err, "service.repository.Create: random error")
	})

	s.Run("should call repository Create method and return error if post transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*domain.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", receiver).Return(nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(errors.New("random error")).Once()

		err := dummyService.CreateReceiver(receiverRequest)
		s.EqualError(err, "slackHelper.PostTransform: random error")
	})
}

func (s *ServiceTestSuite) TestService_GetReceiver() {
	receiverID := uint64(10)
	configurations := make(map[string]interface{})
	configurations["token"] = "key"

	labels := make(map[string]string)
	labels["foo"] = "bar"

	data := make(map[string]interface{})
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
	receiver := &domain.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	s.Run("should call repository Get method and return result in domain's type", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackService: s.slacker}
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(nil).Once()
		s.slacker.On("GetWorkspaceChannels", "key").
			Return([]slack.Channel{
				{ID: "1", Name: "foo"},
			}, nil).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(err)
		s.Equal(dummyReceiver, result)
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should call repository Get method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock}
		newConfigurations := make(map[string]interface{})
		newConfigurations["token"] = "key"
		receiver.Configurations = newConfigurations

		s.repositoryMock.On("Get", receiverID).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should return error if post transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackService: s.slacker}
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "slackHelper.PostTransform: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
		s.slackHelperMock.AssertCalled(s.T(), "PostTransform", receiver)
	})

	s.Run("should return error if getting slack channels failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackService: s.slacker}
		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(nil).Once()
		s.slacker.On("GetWorkspaceChannels", "key").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "could not get channels: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should return error if invalid slack channels", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock, slackService: s.slacker}
		oldjsonMarshal := jsonMarshal
		jsonMarshal = func(interface{}) ([]byte, error) {
			return nil, errors.New("random error")
		}
		defer func() { jsonMarshal = oldjsonMarshal }()

		newConfigurations := make(map[string]interface{})
		newConfigurations["token"] = "key"
		receiver.Configurations = newConfigurations

		s.repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		s.slackHelperMock.On("PostTransform", receiver).
			Return(nil).Once()
		s.slacker.On("GetWorkspaceChannels", "key").
			Return([]slack.Channel{
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
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(map[string]string)
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

	receiver := &domain.Receiver{
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
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*domain.Receiver)
				*r = *receiverResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Update", receiver).Return(nil).Once()

		err := dummyService.UpdateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(receiverResponse, receiverRequest)
		s.repositoryMock.AssertCalled(s.T(), "Update", receiver)
	})

	s.Run("should call repository Create method and return error if transformation failed", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Return(errors.New("random error")).Once()

		err := dummyService.UpdateReceiver(receiverRequest)
		s.EqualError(err, "slackHelper.PreTransform: random error")
	})

	s.Run("should call repository Update method and return error if any", func() {
		dummyService := Service{repository: s.repositoryMock, slackHelper: s.slackHelperMock}
		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*domain.Receiver)
				*r = *receiverResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Update", receiver).
			Return(errors.New("random error")).Once()

		err := dummyService.UpdateReceiver(receiverRequest)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Update", receiver)
	})
}

func (s *ServiceTestSuite) TestService_DeleteReceiver() {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"

	labels := make(map[string]string)
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
