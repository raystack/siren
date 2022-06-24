package subscription_test

import (
	"context"
	"testing"
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/subscription/mocks"
	"github.com/odpf/siren/pkg/cortex"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateSubscription(t *testing.T) {
	receivers := []subscription.ReceiverMetadata{
		{ID: 1, Configuration: map[string]string{"channel_name": "test"}},
	}
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &subscription.Subscription{
		URN:       "foo",
		Namespace: 1,
		Receivers: receivers,
		Match:     match,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	dummyNamespace := &namespace.Namespace{ID: 1, Provider: 1, URN: "dummy"}
	dummyProvider := &provider.Provider{ID: 1, URN: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []receiver.Receiver{
		{ID: 1, Type: "slack", Configurations: map[string]interface{}{"token": "xoxb"}},
		{ID: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{ID: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := new(mocks.CortexClient)
		ctx := context.Background()

		expectedID := uint64(1)
		expectedSubscription := new(subscription.Subscription)
		*expectedSubscription = *input
		expectedSubscription.ID = expectedID
		expectedSubscriptionsInNamespace := []*subscription.Subscription{
			expectedSubscription,
			{
				ID: 2, URN: "bar", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				ID: 3, URN: "baz", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Run(func(ctx context.Context, s *subscription.Subscription) {
			s.ID = expectedID
		}).Once()
		repositoryMock.EXPECT().List(ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), dummyNamespace.URN).
			Run(func(c cortex.AlertManagerConfig, urn string) {
				assert.Len(t, c.Receivers, 3)
				assert.Equal(t, "foo_receiverId_1_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[1].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[2].Receiver)
			}).Return(nil).Once()
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, input)
	})

	t.Run("should return error in subscription creation", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error conflict in subscription creation if repository return error duplicate", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(subscription.ErrDuplicate).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "urn already exist")
	})

	t.Run("should return error in fetching all subscriptions within given namespace", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error in fetching namespace details", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error in fetching provider details", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(&namespace.Namespace{}, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error for unsupported providers", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{input}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).
			Return(&provider.Provider{ID: 1, Type: "prometheus"}, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "subscriptions for provider type 'prometheus' not supported")
	})

	t.Run("should return error in fetching all receivers", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(&namespace.Namespace{}, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(nil, errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error if receiver id not found", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{input}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return([]receiver.Receiver{{ID: 10}}, nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "receiver id 1 does not exist")
	})

	t.Run("should return error if slack channel name not specified in subscription configs", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		dummySubscription := &subscription.Subscription{
			Namespace: input.Namespace,
			Receivers: []subscription.ReceiverMetadata{{ID: 1, Configuration: map[string]string{"id": "1"}}},
		}
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{dummySubscription}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "configuration.channel_name missing from receiver with id 1")
	})

	t.Run("should return error for unsupported receiver type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{input}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return([]receiver.Receiver{{ID: 1, Type: "email"}}, nil).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "subscriptions for receiver type email not supported via Siren inside Cortex")
	})

	t.Run("should return error syncing config with alertmanager", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := &mocks.CortexClient{}
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Create(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{input}, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), "dummy").
			Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "random error")
	})
}

func TestService_GetSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		subsc := &subscription.Subscription{URN: "test", ID: 1, Namespace: 1, Match: make(map[string]string),
			Receivers: []subscription.ReceiverMetadata{{ID: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(subsc, nil).Once()

		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.Equal(t, uint64(1), result.ID)
		assert.Equal(t, "test", result.URN)
		assert.Equal(t, 0, len(result.Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should not return error if subscription doesn't exist", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(nil, nil).Once()
		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("should call repository get method and return error if any", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(nil, errors.New("random error")).Once()
		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.EqualError(t, err, "random error")
		assert.Nil(t, result)
	})

	t.Run("should call repository get method and return error not found if repository return not found error", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(nil, subscription.NotFoundError{}).Once()
		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.EqualError(t, err, "subscription not found")
		assert.Nil(t, result)
	})
}

func TestService_ListSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository list method and return result in domain's type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		subsc := &subscription.Subscription{URN: "test", ID: 1, Namespace: 1, Match: make(map[string]string),
			Receivers: []subscription.ReceiverMetadata{{ID: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.EXPECT().List(ctx).Return([]*subscription.Subscription{subsc}, nil).Once()
		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.ListSubscriptions(context.Background())
		assert.Equal(t, 1, len(result))
		assert.Equal(t, uint64(1), result[0].ID)
		assert.Equal(t, "test", result[0].URN)
		assert.Equal(t, 0, len(result[0].Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should call repository list method and return error if any", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().List(ctx).Return(nil, errors.New("some error")).Once()
		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		result, err := dummyService.ListSubscriptions(context.Background())
		assert.EqualError(t, err, "some error")
		assert.Nil(t, result)
	})
}

func TestService_UpdateSubscription(t *testing.T) {
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &subscription.Subscription{
		ID:        1,
		URN:       "test",
		Namespace: 1,
		Receivers: []subscription.ReceiverMetadata{
			{ID: 1, Configuration: map[string]string{"channel_name": "updated_channel"}},
		},
		Match:     match,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	dummyNamespace := &namespace.Namespace{ID: 1, Provider: 1, URN: "dummy"}
	dummyProvider := &provider.Provider{ID: 1, URN: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []receiver.Receiver{
		{ID: 1, Type: "slack", Configurations: map[string]interface{}{}},
		{ID: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{ID: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := &mocks.CortexClient{}
		ctx := context.Background()

		expectedSubscription := &subscription.Subscription{
			ID:        1,
			URN:       "test",
			Namespace: 1,
			Match:     match,
			Receivers: []subscription.ReceiverMetadata{
				{ID: 1, Configuration: map[string]string{"channel_name": "updated_channel"}},
			},
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		expectedSubscriptionsInNamespace := []*subscription.Subscription{
			{
				ID: input.ID, URN: input.URN, Namespace: input.Namespace,
				Receivers: input.Receivers,
				Match:     input.Match,
			},
			{
				ID: 2, URN: "bar", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				ID: 3, URN: "baz", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Update(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), dummyNamespace.URN).
			Run(func(c cortex.AlertManagerConfig, _a1 string) {
				assert.Len(t, c.Receivers, 3)
				assert.Equal(t, "test_receiverId_1_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[1].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[2].Receiver)
			}).Return(nil).Once()
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, input)
	})

	t.Run("should return error in subscription update", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Update(ctx, input).Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error conflict in subscription update if repository return error duplicate", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Update(ctx, input).Return(subscription.ErrDuplicate).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "urn already exist")
	})

	t.Run("should return error not found in subscription update if repository return not found error", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		ctx := context.Background()

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Update(ctx, input).Return(subscription.NotFoundError{}).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil)

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "subscription not found")
	})

	t.Run("should return error in syncing alertmanager config", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := &mocks.CortexClient{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*subscription.Subscription{
			{
				ID: input.ID, URN: input.URN, Namespace: input.Namespace,
				Receivers: input.Receivers,
				Match:     input.Match,
			},
			{
				ID: 2, URN: "bar", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				ID: 3, URN: "baz", Namespace: input.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Update(ctx, input).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), dummyNamespace.URN).
			Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "random error")
	})
}

func TestService_DeleteSubscription(t *testing.T) {
	match := map[string]string{"foo": "bar"}
	subsc := &subscription.Subscription{
		ID:        1,
		Namespace: 1,
	}
	dummyNamespace := &namespace.Namespace{ID: 1, Provider: 1, URN: "dummy"}
	dummyProvider := &provider.Provider{ID: 1, URN: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []receiver.Receiver{
		{ID: 1, Type: "slack", Configurations: map[string]interface{}{"token": "xoxb"}},
		{ID: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{ID: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository delete method and return result in domain's type", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := &mocks.CortexClient{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*subscription.Subscription{
			{
				ID: 2, URN: "bar", Namespace: subsc.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				ID: 3, URN: "baz", Namespace: subsc.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}
		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(subsc, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Delete(ctx, uint64(1)).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), subsc.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), dummyNamespace.URN).
			Run(func(c cortex.AlertManagerConfig, _ string) {
				assert.Len(t, c.Receivers, 2)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[1].Receiver)
			}).Return(nil).Once()
		repositoryMock.EXPECT().Commit(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("should return error in fetching subscription", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(nil, errors.New("random error")).Once()

		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "random error")
	})

	t.Run("should return error if subscription does not exist", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		ctx := context.Background()

		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(nil, errors.ErrNotFound.WithMsgf("subscription not found")).Once()

		dummyService := subscription.NewService(repositoryMock, nil, nil, nil, nil)

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "subscription not found")
	})

	t.Run("should return error in subscription delete", func(t *testing.T) {
		repositoryMock := new(mocks.SubscriptionRepository)
		providerServiceMock := new(mocks.ProviderService)
		namespaceServiceMock := new(mocks.NamespaceService)
		receiverServiceMock := new(mocks.ReceiverService)
		cortexClientMock := &mocks.CortexClient{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*subscription.Subscription{
			{
				ID: 2, URN: "bar", Namespace: subsc.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				ID: 3, URN: "baz", Namespace: subsc.Namespace,
				Receivers: []subscription.ReceiverMetadata{{ID: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}
		repositoryMock.EXPECT().Get(ctx, uint64(1)).Return(subsc, nil).Once()
		repositoryMock.EXPECT().WithTransaction(ctx).Return(ctx).Once()
		repositoryMock.EXPECT().Delete(ctx, uint64(1)).Return(nil).Once()
		repositoryMock.EXPECT().List(ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), subsc.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(dummyReceivers, nil).Once()
		cortexClientMock.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), dummyNamespace.URN).
			Return(errors.New("random error")).Once()
		repositoryMock.EXPECT().Rollback(ctx).Return(nil).Once()

		dummyService := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "random error")
	})
}
