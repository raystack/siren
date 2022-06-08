package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/core/subscription/alertmanager"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateSubscription(t *testing.T) {
	receivers := []domain.ReceiverMetadata{
		{Id: 1, Configuration: map[string]string{"channel_name": "test"}},
	}
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &domain.Subscription{
		Urn:       "foo",
		Namespace: 1,
		Receivers: receivers,
		Match:     match,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []*domain.Receiver{
		{Id: 1, Type: "slack", Configurations: map[string]interface{}{"token": "xoxb"}},
		{Id: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{Id: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		expectedID := uint64(1)
		expectedSubscription := new(domain.Subscription)
		*expectedSubscription = *input
		expectedSubscription.Id = expectedID
		expectedSubscriptionsInNamespace := []*domain.Subscription{
			expectedSubscription,
			{
				Id: 2, Urn: "bar", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				Id: 3, Urn: "baz", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).
			Run(func(args mock.Arguments) {
				s := args.Get(1).(*domain.Subscription)
				s.Id = expectedID
			}).Once()
		repositoryMock.On("List", ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), dummyNamespace.Urn).
			Run(func(args mock.Arguments) {
				c := args.Get(0).(alertmanager.AMConfig)
				assert.Len(t, c.Receivers, 3)
				assert.Equal(t, "foo_receiverId_1_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[1].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[2].Receiver)
			}).Return(nil).Once()
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, input)
	})

	t.Run("should return error in subscription creation", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.repository.Create: random error")
	})

	t.Run("should return error in fetching all subscriptions within given namespace", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.getAllSubscriptionsWithinNamespace: s.repository.List: random error")
	})

	t.Run("should return error in fetching namespace details", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.getProviderAndNamespaceInfoFromNamespaceId: failed to get namespace details: random error")
	})

	t.Run("should return error in fetching provider details", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(&domain.Namespace{}, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.getProviderAndNamespaceInfoFromNamespaceId: failed to get provider details: random error")
	})

	t.Run("should return error for unsupported providers", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{input}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).
			Return(&domain.Provider{Id: 1, Type: "prometheus"}, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: subscriptions for provider type 'prometheus' not supported")
	})

	t.Run("should return error in fetching all receivers", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(&domain.Namespace{}, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(&domain.Provider{}, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(nil, errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.addReceiversConfiguration: failed to get receivers: random error")
	})

	t.Run("should return error if receiver id not found", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{input}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return([]*domain.Receiver{{Id: 10}}, nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.addReceiversConfiguration: receiver id 1 does not exist")
	})

	t.Run("should return error if slack channel name not specified in subscription configs", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		dummySubscription := &domain.Subscription{
			Namespace: input.Namespace,
			Receivers: []domain.ReceiverMetadata{{Id: 1, Configuration: map[string]string{"id": "1"}}},
		}
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{dummySubscription}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.addReceiversConfiguration: configuration.channel_name missing from receiver with id 1")
	})

	t.Run("should return error for unsupported receiver type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{input}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return([]*domain.Receiver{{Id: 1, Type: "email"}}, nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.addReceiversConfiguration: subscriptions for receiver type email not supported via Siren inside Cortex")
	})

	t.Run("should return error in alertmanager client initialization", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{input}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return nil, errors.New("random error")
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: alertmanagerClientCreator: : random error")
	})

	t.Run("should return error syncing config with alertmanager", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Create", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{input}, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", mock.AnythingOfType("uint64")).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").
			Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}
		err := dummyService.CreateSubscription(context.Background(), input)

		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.amClient.SyncConfig: random error")
	})
}

func TestService_GetSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository get method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		subscription := &domain.Subscription{Urn: "test", Id: 1, Namespace: 1, Match: make(map[string]string),
			Receivers: []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("Get", ctx, uint64(1)).Return(subscription, nil).Once()

		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.Equal(t, uint64(1), result.Id)
		assert.Equal(t, "test", result.Urn)
		assert.Equal(t, 0, len(result.Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should not return error if subscription doesn't exist", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		repositoryMock.On("Get", ctx, uint64(1)).Return(nil, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("should call repository get method and return error if any", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		repositoryMock.On("Get", ctx, uint64(1)).Return(nil, errors.New("random error")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		result, err := dummyService.GetSubscription(context.Background(), 1)
		assert.EqualError(t, err, "s.repository.Get: random error")
		assert.Nil(t, result)
	})
}

func TestService_ListSubscription(t *testing.T) {
	timeNow := time.Now()

	t.Run("should call repository list method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		subscription := &domain.Subscription{Urn: "test", Id: 1, Namespace: 1, Match: make(map[string]string),
			Receivers: []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
			CreatedAt: timeNow, UpdatedAt: timeNow}
		repositoryMock.On("List", ctx).Return([]*domain.Subscription{subscription}, nil).Once()
		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		result, err := dummyService.ListSubscriptions(context.Background())
		assert.Equal(t, 1, len(result))
		assert.Equal(t, uint64(1), result[0].Id)
		assert.Equal(t, "test", result[0].Urn)
		assert.Equal(t, 0, len(result[0].Receivers[0].Configuration))
		assert.Nil(t, err)
	})

	t.Run("should call repository list method and return error if any", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		repositoryMock.On("List", ctx).Return(nil, errors.New("abcd")).Once()
		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		result, err := dummyService.ListSubscriptions(context.Background())
		assert.EqualError(t, err, "s.repository.List: abcd")
		assert.Nil(t, result)
	})
}

func TestService_UpdateSubscription(t *testing.T) {
	match := make(map[string]string)
	match["foo"] = "bar"
	timeNow := time.Now()
	input := &domain.Subscription{
		Id:        1,
		Urn:       "test",
		Namespace: 1,
		Receivers: []domain.ReceiverMetadata{
			{Id: 1, Configuration: map[string]string{"channel_name": "updated_channel"}},
		},
		Match:     match,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []*domain.Receiver{
		{Id: 1, Type: "slack", Configurations: map[string]interface{}{}},
		{Id: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{Id: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		expectedSubscription := &domain.Subscription{
			Id:        1,
			Urn:       "test",
			Namespace: 1,
			Match:     match,
			Receivers: []domain.ReceiverMetadata{
				{Id: 1, Configuration: map[string]string{"channel_name": "updated_channel"}},
			},
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		expectedSubscriptionsInNamespace := []*domain.Subscription{
			{
				Id: input.Id, Urn: input.Urn, Namespace: input.Namespace,
				Receivers: input.Receivers,
				Match:     input.Match,
			},
			{
				Id: 2, Urn: "bar", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				Id: 3, Urn: "baz", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Update", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), dummyNamespace.Urn).
			Run(func(args mock.Arguments) {
				c := args.Get(0).(alertmanager.AMConfig)
				assert.Len(t, c.Receivers, 3)
				assert.Equal(t, "test_receiverId_1_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[1].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[2].Receiver)
			}).Return(nil).Once()
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, expectedSubscription, input)
	})

	t.Run("should return error in subscription update", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		ctx := context.Background()

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Update", ctx, input).Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, nil}

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "s.repository.Update: random error")
	})

	t.Run("should return error in syncing alertmanager config", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*domain.Subscription{
			{
				Id: input.Id, Urn: input.Urn, Namespace: input.Namespace,
				Receivers: input.Receivers,
				Match:     input.Match,
			},
			{
				Id: 2, Urn: "bar", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				Id: 3, Urn: "baz", Namespace: input.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}

		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Update", ctx, input).Return(nil).Once()
		repositoryMock.On("List", ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.On("GetNamespace", input.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), dummyNamespace.Urn).
			Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}

		err := dummyService.UpdateSubscription(context.Background(), input)
		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.amClient.SyncConfig: random error")
	})
}

func TestService_DeleteSubscription(t *testing.T) {
	match := map[string]string{"foo": "bar"}
	subscription := &domain.Subscription{
		Id:        1,
		Namespace: 1,
	}
	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}
	dummyReceivers := []*domain.Receiver{
		{Id: 1, Type: "slack", Configurations: map[string]interface{}{"token": "xoxb"}},
		{Id: 2, Type: "pagerduty", Configurations: map[string]interface{}{"service_key": "abcd"}},
		{Id: 3, Type: "http", Configurations: map[string]interface{}{"url": "http://localhost:3000"}},
	}

	t.Run("should call repository delete method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*domain.Subscription{
			{
				Id: 2, Urn: "bar", Namespace: subscription.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				Id: 3, Urn: "baz", Namespace: subscription.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}
		repositoryMock.On("Get", ctx, uint64(1)).Return(subscription, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Delete", ctx, uint64(1)).Return(nil).Once()
		repositoryMock.On("List", ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.On("GetNamespace", subscription.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), dummyNamespace.Urn).
			Run(func(args mock.Arguments) {
				c := args.Get(0).(alertmanager.AMConfig)
				assert.Len(t, c.Receivers, 2)
				assert.Equal(t, "bar_receiverId_2_idx_0", c.Receivers[0].Receiver)
				assert.Equal(t, "baz_receiverId_3_idx_0", c.Receivers[1].Receiver)
			}).Return(nil).Once()
		repositoryMock.On("Commit", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.Nil(t, err)
	})

	t.Run("should return error in fetching subscription", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		repositoryMock.On("Get", ctx, uint64(1)).Return(nil, errors.New("random error")).Once()

		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "s.repository.Get: random error")
	})

	t.Run("should return error if subscription does not exist", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		ctx := context.Background()

		repositoryMock.On("Get", ctx, uint64(1)).Return(nil, nil).Once()

		dummyService := Service{repositoryMock, nil, nil, nil, nil}

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "subscription not found")
	})

	t.Run("should return error in subscription delete", func(t *testing.T) {
		repositoryMock := &SubscriptionRepositoryMock{}
		providerServiceMock := &mocks.ProviderService{}
		namespaceServiceMock := &mocks.NamespaceService{}
		receiverServiceMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		ctx := context.Background()

		expectedSubscriptionsInNamespace := []*domain.Subscription{
			{
				Id: 2, Urn: "bar", Namespace: subscription.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 2, Configuration: make(map[string]string)}},
				Match:     match,
			},
			{
				Id: 3, Urn: "baz", Namespace: subscription.Namespace,
				Receivers: []domain.ReceiverMetadata{{Id: 3, Configuration: make(map[string]string)}},
				Match:     match,
			},
		}
		repositoryMock.On("Get", ctx, uint64(1)).Return(subscription, nil).Once()
		repositoryMock.On("WithTransaction", ctx).Return(ctx).Once()
		repositoryMock.On("Delete", ctx, uint64(1)).Return(nil).Once()
		repositoryMock.On("List", ctx).Return(expectedSubscriptionsInNamespace, nil).Once()
		namespaceServiceMock.On("GetNamespace", subscription.Namespace).Return(dummyNamespace, nil).Once()
		providerServiceMock.On("GetProvider", dummyNamespace.Provider).Return(dummyProvider, nil).Once()
		receiverServiceMock.On("ListReceivers").Return(dummyReceivers, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), dummyNamespace.Urn).
			Return(errors.New("random error")).Once()
		repositoryMock.On("Rollback", ctx).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		dummyService := Service{repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, amClientMock}

		err := dummyService.DeleteSubscription(context.Background(), 1)
		assert.EqualError(t, err, "s.syncInUpstreamCurrentSubscriptionsOfNamespace: s.amClient.SyncConfig: random error")
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
