package receiver_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	"github.com/odpf/siren/plugins/receivers/slack"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	jsonMarshal = json.Marshal
)

type ServiceTestSuite struct {
	suite.Suite
	slackHelperMock *receiver.MockSlackHelper
	repositoryMock  *mocks.ReceiverRepository
	slacker         *mocks.SlackService
}

func (s *ServiceTestSuite) SetupTest() {
	s.slackHelperMock = &receiver.MockSlackHelper{}
	s.repositoryMock = &mocks.ReceiverRepository{}
	s.slacker = &mocks.SlackService{}
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestService_ListReceivers() {
	configurations := make(map[string]interface{})
	labels := make(map[string]string)
	labels["foo"] = "bar"

	s.Run("should call repository List method and return result in domain's type", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		dummyReceivers := []*receiver.Receiver{
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
		receivers := []*receiver.Receiver{
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
		dummyService, err := receiver.NewService(s.repositoryMock, nil, s.slacker)
		s.NoError(err)

		configurations["token"] = "key"
		s.repositoryMock.On("List").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.ListReceivers()
		s.Nil(result)
		s.EqualError(err, "service.repository.List: random error")
		s.repositoryMock.AssertCalled(s.T(), "List")
	})

	s.Run("should call repository List method and return error if post slack transformation failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		receivers := []*receiver.Receiver{
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
	receiverRequest := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	transfromResponse := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: tokenConfigurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	rcv := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: tokenConfigurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	s.Run("should call repository Create method and return result in domain's type", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*receiver.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", rcv).Return(nil).Once()
		s.slackHelperMock.On("PostTransform", rcv).
			Return(nil).Once()

		err = dummyService.CreateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(transfromResponse, receiverRequest)
		s.repositoryMock.AssertCalled(s.T(), "Create", rcv)
	})

	s.Run("should call repository Create method and return error if pre transformation failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Return(errors.New("random error")).Once()

		err = dummyService.CreateReceiver(receiverRequest)
		s.EqualError(err, "slackHelper.PreTransform: random error")
	})

	s.Run("should call repository Create method and return error if any", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*receiver.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", rcv).
			Return(errors.New("random error")).Once()

		err = dummyService.CreateReceiver(receiverRequest)
		s.EqualError(err, "service.repository.Create: random error")
	})

	s.Run("should call repository Create method and return error if post transformation failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*receiver.Receiver)
				*r = *transfromResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Create", rcv).Return(nil).Once()
		s.slackHelperMock.On("PostTransform", rcv).
			Return(errors.New("random error")).Once()

		err = dummyService.CreateReceiver(receiverRequest)
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
	dummyReceiver := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		Data:           data,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	rcv := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	s.Run("should call repository Get method and return result in domain's type", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Get", receiverID).Return(rcv, nil).Once()
		s.slackHelperMock.On("PostTransform", rcv).
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
		dummyService, err := receiver.NewService(s.repositoryMock, nil, s.slacker)
		s.NoError(err)

		newConfigurations := make(map[string]interface{})
		newConfigurations["token"] = "key"
		rcv.Configurations = newConfigurations

		s.repositoryMock.On("Get", receiverID).
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	s.Run("should return error if post transformation failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Get", receiverID).Return(rcv, nil).Once()
		s.slackHelperMock.On("PostTransform", rcv).
			Return(errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "slackHelper.PostTransform: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
		s.slackHelperMock.AssertCalled(s.T(), "PostTransform", rcv)
	})

	s.Run("should return error if getting slack channels failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Get", receiverID).Return(rcv, nil).Once()
		s.slackHelperMock.On("PostTransform", rcv).
			Return(nil).Once()
		s.slacker.On("GetWorkspaceChannels", "key").
			Return(nil, errors.New("random error")).Once()

		result, err := dummyService.GetReceiver(receiverID)
		s.Nil(result)
		s.EqualError(err, "could not get channels: random error")
		s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	})

	//TODO confusing logic test
	// s.Run("should return error if invalid slack channels", func() {
	// 	dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
	// 	s.NoError(err)

	// 	oldjsonMarshal := jsonMarshal
	// 	jsonMarshal = func(interface{}) ([]byte, error) {
	// 		return nil, errors.New("random error")
	// 	}
	// 	defer func() { jsonMarshal = oldjsonMarshal }()

	// 	newConfigurations := make(map[string]interface{})
	// 	newConfigurations["token"] = "key"
	// 	rcv.Configurations = newConfigurations

	// 	s.repositoryMock.On("Get", receiverID).Return(rcv, nil).Once()
	// 	s.slackHelperMock.On("PostTransform", rcv).
	// 		Return(nil).Once()
	// 	s.slacker.On("GetWorkspaceChannels", "key").
	// 		Return([]slack.Channel{
	// 			{ID: "1", Name: string([]byte{0xff})},
	// 		}, nil).Once()

	// 	result, err := dummyService.GetReceiver(receiverID)
	// 	s.Nil(result)
	// 	s.EqualError(err, "invalid channels: random error")
	// 	s.repositoryMock.AssertCalled(s.T(), "Get", receiverID)
	// })
}

func (s *ServiceTestSuite) TestService_UpdateReceiver() {
	timenow := time.Now()
	configurations := make(map[string]interface{})
	configurations["client_id"] = "foo"
	configurations["client_secret"] = "bar"
	configurations["auth_code"] = "foo"

	labels := make(map[string]string)
	labels["foo"] = "bar"
	receiverRequest := &receiver.Receiver{
		Id:             10,
		Name:           "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	rcv := &receiver.Receiver{
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

	receiverResponse := &receiver.Receiver{
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
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*receiver.Receiver)
				*r = *receiverResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Update", rcv).Return(nil).Once()

		err = dummyService.UpdateReceiver(receiverRequest)
		s.Nil(err)
		s.Equal(receiverResponse, receiverRequest)
		s.repositoryMock.AssertCalled(s.T(), "Update", rcv)
	})

	s.Run("should call repository Create method and return error if transformation failed", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Return(errors.New("random error")).Once()

		err = dummyService.UpdateReceiver(receiverRequest)
		s.EqualError(err, "slackHelper.PreTransform: random error")
	})

	s.Run("should call repository Update method and return error if any", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, s.slackHelperMock, s.slacker)
		s.NoError(err)

		s.slackHelperMock.On("PreTransform", receiverRequest).
			Run(func(args mock.Arguments) {
				r := args.Get(0).(*receiver.Receiver)
				*r = *receiverResponse
			}).Return(nil).Once()
		s.repositoryMock.On("Update", rcv).
			Return(errors.New("random error")).Once()

		err = dummyService.UpdateReceiver(receiverRequest)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Update", rcv)
	})
}

func (s *ServiceTestSuite) TestService_DeleteReceiver() {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"

	labels := make(map[string]string)
	labels["foo"] = "bar"
	receiverID := uint64(10)

	s.Run("should call repository Delete method and return nil if no error", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, nil, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Delete", receiverID).Return(nil).Once()

		err = dummyService.DeleteReceiver(receiverID)
		s.Nil(err)
		s.repositoryMock.AssertCalled(s.T(), "Delete", receiverID)
	})

	s.Run("should call repository Delete method and return error if any", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, nil, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Delete", receiverID).
			Return(errors.New("random error")).Once()

		err = dummyService.DeleteReceiver(receiverID)
		s.EqualError(err, "random error")
		s.repositoryMock.AssertCalled(s.T(), "Delete", receiverID)
	})
}

func (s *ServiceTestSuite) TestService_Migrate() {
	s.Run("should call repository Migrate method and return result", func() {
		dummyService, err := receiver.NewService(s.repositoryMock, nil, s.slacker)
		s.NoError(err)

		s.repositoryMock.On("Migrate").Return(nil).Once()

		err = dummyService.Migrate()
		s.Nil(err)
		s.repositoryMock.AssertCalled(s.T(), "Migrate")
	})
}
