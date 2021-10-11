package receiver

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestListReceivers(t *testing.T) {
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"

	t.Run("should call repository List method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		dummyReceivers := []*domain.Receiver{
			{
				Id:             10,
				Urn:            "foo",
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
				Urn:            "foo",
				Type:           "slack",
				Labels:         labels,
				Configurations: configurations,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		repositoryMock.On("List").Return(receivers, nil).Once()
		result, err := dummyService.ListReceivers()
		assert.Nil(t, err)
		assert.Equal(t, len(dummyReceivers), len(result))
		assert.Equal(t, dummyReceivers[0].Urn, result[0].Urn)
		repositoryMock.AssertCalled(t, "List")
	})

	t.Run("should call repository List method and return error if any", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("List").
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.ListReceivers()
		assert.Nil(t, result)
		assert.EqualError(t, err, "service.repository.List: random error")
		repositoryMock.AssertCalled(t, "List")
	})
}

func TestCreateReceiver(t *testing.T) {
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyReceiver := &domain.Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	receiver := &Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	t.Run("should call repository Create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Create", receiver).Return(receiver, nil).Once()
		result, err := dummyService.CreateReceiver(dummyReceiver)
		assert.Nil(t, err)
		assert.Equal(t, dummyReceiver, result)
		repositoryMock.AssertCalled(t, "Create", receiver)
	})

	t.Run("should call repository Create method and return error if any", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Create", receiver).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.CreateReceiver(dummyReceiver)
		assert.Nil(t, result)
		assert.EqualError(t, err, "service.repository.Create: random error")
		repositoryMock.AssertCalled(t, "Create", receiver)
	})
}

func TestGetReceiver(t *testing.T) {
	receiverID := uint64(10)
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	timenow := time.Now()
	dummyReceiver := &domain.Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	receiver := &Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	t.Run("should call repository Get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", receiverID).Return(receiver, nil).Once()
		result, err := dummyService.GetReceiver(receiverID)
		assert.Nil(t, err)
		assert.Equal(t, dummyReceiver, result)
		repositoryMock.AssertCalled(t, "Get", receiverID)
	})

	t.Run("should call repository Get method and return error if any", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Get", receiverID).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.GetReceiver(receiverID)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Get", receiverID)
	})
}

func TestUpdateReceiver(t *testing.T) {
	timenow := time.Now()
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	dummyReceiver := &domain.Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}
	receiver := &Receiver{
		Id:             10,
		Urn:            "foo",
		Type:           "slack",
		Labels:         labels,
		Configurations: configurations,
		CreatedAt:      timenow,
		UpdatedAt:      timenow,
	}

	t.Run("should call repository Update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Update", receiver).Return(receiver, nil).Once()
		result, err := dummyService.UpdateReceiver(dummyReceiver)
		assert.Nil(t, err)
		assert.Equal(t, dummyReceiver, result)
		repositoryMock.AssertCalled(t, "Update", receiver)
	})

	t.Run("should call repository Update method and return error if any", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Update", receiver).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.UpdateReceiver(dummyReceiver)
		assert.Nil(t, result)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Update", receiver)
	})
}

func TestDeleteReceiver(t *testing.T) {
	configurations := make(StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(StringStringMap)
	labels["foo"] = "bar"
	receiverID := uint64(10)

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", receiverID).Return(nil).Once()
		err := dummyService.DeleteReceiver(receiverID)
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Delete", receiverID)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Delete", receiverID).
			Return(errors.New("random error")).Once()
		err := dummyService.DeleteReceiver(receiverID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Delete", receiverID)
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &MockReceiverRepository{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
