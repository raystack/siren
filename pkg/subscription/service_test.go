package subscription

import (
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_CreateSubscription(t *testing.T) {
	receivers := []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &domain.Subscription{Id: 1, Urn: "test", Namespace: 1, Receivers: receivers, Match: match, CreatedAt: timeNow, UpdatedAt: timeNow}

	t.Run("should call repository create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: match,
			Receiver: []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}, CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Create", subscription, nil, nil, nil).Return(subscription, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.CreateSubscription(input)
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, uint64(1), result.Receivers[0].Id)
		assert.Equal(t, 0, len(result.Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should return error in subscription creation", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: match,
			Receiver: []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}, CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Create", subscription, nil, nil, nil).Return(nil, errors.New("random error")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.CreateSubscription(input)
		assert.EqualError(t, err, "s.repository.Create: random error")
		assert.Nil(t, result)
	})
}

func TestService_GetSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: make(map[string]string),
			Receiver:  []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Get", uint64(1)).Return(subscription, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.GetSubscription(1)
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, "test", result.Urn)
		assert.Equal(t, 0, len(result.Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should not return error if subscription doesn't exist", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		repositoryMock.On("Get", uint64(1)).Return(nil, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.GetSubscription(1)
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("should call repository get method and return error if any", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		repositoryMock.On("Get", uint64(1)).Return(nil, errors.New("random error")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.GetSubscription(1)
		assert.EqualError(t, err, "s.repository.Get: random error")
		assert.Nil(t, result)
	})
}

func TestService_ListSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository list method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: make(map[string]string),
			Receiver:  []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("List").Return([]*Subscription{subscription}, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.ListSubscriptions()
		assert.Equal(t, 1, len(result))
		assert.Equal(t, uint64(1), result[0].Id)
		assert.Equal(t, "test", result[0].Urn)
		assert.Equal(t, 0, len(result[0].Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should call repository list method and return error if any", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		repositoryMock.On("List").Return(nil, errors.New("abcd")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.ListSubscriptions()
		assert.EqualError(t, err, "s.repository.List: abcd")
		assert.Nil(t, result)
	})
}

func TestService_UpdateSubscription(t *testing.T) {
	receivers := []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &domain.Subscription{Id: 1, Urn: "test", Namespace: 1, Receivers: receivers, Match: match, CreatedAt: timeNow, UpdatedAt: timeNow}

	t.Run("should call repository update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: match,
			Receiver: []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}, CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Update", subscription, nil, nil, nil).Return(subscription, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.UpdateSubscription(input)
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, uint64(1), result.Receivers[0].Id)
		assert.Equal(t, 0, len(result.Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should return error in subscription update", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		subscription := &Subscription{Urn: "test", Id: 1, NamespaceId: 1, Match: match,
			Receiver: []ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}}, CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Update", subscription, nil, nil, nil).Return(nil, errors.New("random error")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		result, err := dummyService.UpdateSubscription(input)
		assert.EqualError(t, err, "s.repository.Update: random error")
		assert.Nil(t, result)
	})
}

func TestService_DeleteSubscription(t *testing.T) {
	t.Run("should call repository delete method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		repositoryMock.On("Delete", uint64(1), nil, nil, nil).Return(nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		err := dummyService.DeleteSubscription(1)
		assert.Nil(t, err)
	})

	t.Run("should return error in subscription delete", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		repositoryMock.On("Delete", uint64(1), nil, nil, nil).Return(errors.New("random error")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil}

		err := dummyService.DeleteSubscription(1)
		assert.EqualError(t, err, "random error")
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		dummyService := Service{repository: repositoryMock}
		repositoryMock.On("Migrate").Return(nil).Once()
		err := dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}
