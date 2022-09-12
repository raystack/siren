package subscription_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/subscription/mocks"
	"github.com/stretchr/testify/mock"
)

func TestService_List(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.SubscriptionRepository)
		ErrString   string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return nil error if subscription list repository return no error",
				Setup: func(sr *mocks.SubscriptionRepository) {
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{}, nil)
				},
			},
			{
				Description: "should return error if subscription list repository return error",
				Setup: func(sr *mocks.SubscriptionRepository) {
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.SubscriptionRepository)
			)
			svc := subscription.NewService(repositoryMock, nil, nil, nil, nil)

			tc.Setup(repositoryMock)

			_, err := svc.List(ctx, subscription.Filter{})
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.SubscriptionRepository)
		ErrString   string
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return nil error if subscription get repository return no error",
				Setup: func(sr *mocks.SubscriptionRepository) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
				},
			},
			{
				Description: "should return error if subscription get repository return error",
				Setup: func(sr *mocks.SubscriptionRepository) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error not found if subscription get repository return not found",
				Setup: func(sr *mocks.SubscriptionRepository) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, subscription.NotFoundError{})
				},
				ErrString: "subscription not found",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.SubscriptionRepository)
			)
			svc := subscription.NewService(repositoryMock, nil, nil, nil, nil)

			tc.Setup(repositoryMock)

			_, err := svc.Get(ctx, 100)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestService_Create(t *testing.T) {
	type testCase struct {
		Description  string
		Subscription *subscription.Subscription
		Setup        func(*mocks.SubscriptionRepository, *mocks.NamespaceService, *mocks.ProviderService, *mocks.ReceiverService, *mocks.CortexClient)
		ErrString    string
	}

	var (
		ctx                       = context.TODO()
		mockSyncToUpstreamSuccess = func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {

			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
				{
					ID: 1,
				},
			}, nil)
			rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
			cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), mock.AnythingOfType("string")).Return(nil)
		}
		testCases = []testCase{
			{
				Description: "should return error if namespace service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if provider service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error conflict if create subscription return error duplicate",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(subscription.ErrDuplicate)
					sr.EXPECT().Rollback(ctx, subscription.ErrDuplicate).Return(nil)
				},
				ErrString: "urn already exist",
			},
			{
				Description: "should return error not found if create subscription return error relation",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(subscription.ErrRelation)
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "namespace id does not exist",
			},
			{
				Description: "should return error if create subscription return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return no error if create subscription return no error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return error if transaction commit return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(errors.New("commit error"))
				},
				ErrString: "commit error",
			},
			{
				Description: "should return error if create repository return error and rollback return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback success",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback failed",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock       = new(mocks.SubscriptionRepository)
				namespaceServiceMock = new(mocks.NamespaceService)
				providerServiceMock  = new(mocks.ProviderService)
				receiverServiceMock  = new(mocks.ReceiverService)
				cortexClientMock     = new(mocks.CortexClient)
			)
			svc := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)
			tc.Setup(repositoryMock, namespaceServiceMock, providerServiceMock, receiverServiceMock, cortexClientMock)

			err := svc.Create(ctx, &subscription.Subscription{})
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			providerServiceMock.AssertExpectations(t)
			receiverServiceMock.AssertExpectations(t)
			cortexClientMock.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	type testCase struct {
		Description  string
		Subscription *subscription.Subscription
		Setup        func(*mocks.SubscriptionRepository, *mocks.NamespaceService, *mocks.ProviderService, *mocks.ReceiverService, *mocks.CortexClient)
		ErrString    string
	}
	var (
		ctx                       = context.TODO()
		mockSyncToUpstreamSuccess = func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
				{
					ID: 1,
				},
			}, nil)
			rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
			cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), mock.AnythingOfType("string")).Return(nil)
		}
		testCases = []testCase{
			{
				Description: "should return error if namespace service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if provider service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error conflict if update subscription return error duplicate",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(subscription.ErrDuplicate)
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "urn already exist",
			},
			{
				Description: "should return error not found if update subscription return error relation",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(subscription.ErrRelation)
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "namespace id does not exist",
			},
			{
				Description: "should return error not found if update subscription return not found error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(subscription.NotFoundError{})
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "subscription not found",
			},
			{
				Description: "should return error if update subscription return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if update subscription return error and rollback return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback success",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback failed",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
			{
				Description: "should return no error if update subscription return no error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return error if transaction commit return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*subscription.Subscription")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(errors.New("commit error"))
				},
				ErrString: "commit error",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock       = new(mocks.SubscriptionRepository)
				namespaceServiceMock = new(mocks.NamespaceService)
				providerServiceMock  = new(mocks.ProviderService)
				receiverServiceMock  = new(mocks.ReceiverService)
				cortexClientMock     = new(mocks.CortexClient)
			)
			svc := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)
			tc.Setup(repositoryMock, namespaceServiceMock, providerServiceMock, receiverServiceMock, cortexClientMock)

			err := svc.Update(ctx, &subscription.Subscription{})
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			providerServiceMock.AssertExpectations(t)
			receiverServiceMock.AssertExpectations(t)
			cortexClientMock.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	type testCase struct {
		Description  string
		Subscription *subscription.Subscription
		Setup        func(*mocks.SubscriptionRepository, *mocks.NamespaceService, *mocks.ProviderService, *mocks.ReceiverService, *mocks.CortexClient)
		ErrString    string
	}
	var (
		ctx                       = context.TODO()
		mockSyncToUpstreamSuccess = func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
				{
					URN:       "alert-history-odpf",
					Namespace: 2,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
				{
					URN:       "odpf-data-warning",
					Namespace: 1,
					Receivers: []subscription.Receiver{
						{
							ID: 1,
							Configuration: map[string]string{
								"channel_name": "odpf-data",
							},
						},
					},
					Match: map[string]string{
						"environment": "integration",
						"team":        "odpf-data",
					},
				},
			}, nil)
			rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
				{
					ID: 1,
				},
			}, nil)
			rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
			cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), mock.AnythingOfType("string")).Return(nil)
		}
		testCases = []testCase{
			{
				Description: "should return error if get subscription repository return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if namespace service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if provider service get return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if delete subscription return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if delete subscription return error and rollback return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback success",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(nil)
				},
				ErrString: "some error",
			},
			{
				Description: "should return error if sync to upstream return error and rollback failed",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil)
					sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
					sr.EXPECT().Rollback(ctx, mock.Anything).Return(errors.New("some rollback error"))
				},
				ErrString: "some rollback error",
			},
			{
				Description: "should return no error if delete subscription return no error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(nil)
				},
			},
			{
				Description: "should return error if transaction commit return error",
				Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
					sr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&subscription.Subscription{}, nil)
					ns.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&namespace.Namespace{}, nil)
					ps.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&provider.Provider{
						Type: provider.TypeCortex,
					}, nil)

					sr.EXPECT().WithTransaction(ctx).Return(ctx)
					sr.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil)
					mockSyncToUpstreamSuccess(sr, ns, ps, rs, cc)
					sr.EXPECT().Commit(ctx).Return(errors.New("commit error"))
				},
				ErrString: "commit error",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock       = new(mocks.SubscriptionRepository)
				namespaceServiceMock = new(mocks.NamespaceService)
				providerServiceMock  = new(mocks.ProviderService)
				receiverServiceMock  = new(mocks.ReceiverService)
				cortexClientMock     = new(mocks.CortexClient)
			)
			svc := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)
			tc.Setup(repositoryMock, namespaceServiceMock, providerServiceMock, receiverServiceMock, cortexClientMock)

			err := svc.Delete(ctx, 100)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			providerServiceMock.AssertExpectations(t)
			receiverServiceMock.AssertExpectations(t)
			cortexClientMock.AssertExpectations(t)
		})
	}
}

func TestCreateReceiversMap(t *testing.T) {

	type testCase struct {
		Description          string
		Subscriptions        []subscription.Subscription
		Setup                func(*mocks.ReceiverService)
		ExpectedReceiversMap map[uint64]*receiver.Receiver
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should return error if subscription does not have receivers",
			Setup:       func(rs *mocks.ReceiverService) {},
			ErrString:   "no receivers found in subscription",
		},
		{
			Description: "should return error if receiver service list return error",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description: "should return error empty receivers if at least one receiver ids in map don't exist",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11},
				}, nil)
			},
			ErrString: "receiver id [12] don't exist",
		},
		{
			Description: "should return receivers map if all receiver ids exist",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
					},
				},
			},
			ExpectedReceiversMap: map[uint64]*receiver.Receiver{
				1:  {ID: 1},
				2:  {ID: 2},
				5:  {ID: 5},
				7:  {ID: 7},
				9:  {ID: 9},
				11: {ID: 11},
				12: {ID: 12},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				receiverServiceMock = new(mocks.ReceiverService)
			)

			tc.Setup(receiverServiceMock)

			got, err := subscription.CreateReceiversMap(context.Background(), receiverServiceMock, tc.Subscriptions)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedReceiversMap) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedReceiversMap)
			}
		})
	}
}

func TestAssignReceivers(t *testing.T) {

	type testCase struct {
		Description           string
		Subscriptions         []subscription.Subscription
		ReceiversMap          map[uint64]*receiver.Receiver
		Setup                 func(*mocks.ReceiverService)
		ExpectedSubscriptions []subscription.Subscription
		ErrString             string
	}

	var inputSubscriptions = []subscription.Subscription{
		{
			Receivers: []subscription.Receiver{
				{
					ID: 1,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
				{
					ID: 2,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
			},
		},
		{
			Receivers: []subscription.Receiver{
				{
					ID: 3,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
				{
					ID: 4,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
			},
		},
	}

	var testCases = []testCase{
		{
			Description:   "should return error if can't found receiver id from map",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: nil,
				2: nil,
				3: nil,
				4: nil,
			},
			Setup: func(rs *mocks.ReceiverService) {
			},
			ErrString: "receiver id 1 not found",
		},
		{
			Description:   "should return error if get subscription config return error",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: {ID: 1, Type: receiver.TypeHTTP},
				2: {ID: 2, Type: receiver.TypePagerDuty},
				3: {ID: 3, Type: receiver.TypeSlack},
				4: {ID: 4, Type: receiver.TypeSlack},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description:   "should assign receivers to subscription if assigning receivers return no error",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: {ID: 1, Type: receiver.TypeHTTP},
				2: {ID: 2, Type: receiver.TypePagerDuty},
				3: {ID: 3, Type: receiver.TypeSlack},
				4: {ID: 4, Type: receiver.TypeSlack},
			},
			ExpectedSubscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{
							ID:   1,
							Type: receiver.TypeHTTP,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
						{
							ID:   2,
							Type: receiver.TypePagerDuty,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
					},
				},
				{
					Receivers: []subscription.Receiver{
						{
							ID:   3,
							Type: receiver.TypeSlack,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
						{
							ID:   4,
							Type: receiver.TypeSlack,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{
					"token":  "abcabc",
					"newkey": "newvalue",
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				receiverServiceMock = new(mocks.ReceiverService)
			)

			tc.Setup(receiverServiceMock)

			got, err := subscription.AssignReceivers(receiverServiceMock, tc.ReceiversMap, tc.Subscriptions)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedSubscriptions) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedSubscriptions)
			}
		})
	}
}

func TestSyncToUpstream(t *testing.T) {

	type testCase struct {
		Description string
		Setup       func(*mocks.SubscriptionRepository, *mocks.NamespaceService, *mocks.ProviderService, *mocks.ReceiverService, *mocks.CortexClient)
		Namespace   *namespace.Namespace
		Provider    *provider.Provider
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should return error if list subscriptions in namespace return error",
			Namespace: &namespace.Namespace{
				ID: 111,
			},
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description: "should return error if create receivers map return error",
			Namespace: &namespace.Namespace{
				ID: 111,
			},
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{}, nil)
			},
			ErrString: "no receivers found in subscription",
		},
		{
			Description: "should return error if assign receivers return error",
			Namespace: &namespace.Namespace{
				ID: 111,
			},
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
					{
						URN:       "alert-history-odpf",
						Namespace: 2,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
					{
						URN:       "odpf-data-warning",
						Namespace: 1,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
								Configuration: map[string]string{
									"channel_name": "odpf-data",
								},
							},
						},
						Match: map[string]string{
							"environment": "integration",
							"team":        "odpf-data",
						},
					},
				}, nil)
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{
						ID: 1,
					},
				}, nil)
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description: "should return error if provider type is unknown",
			Namespace: &namespace.Namespace{
				ID: 111,
			},
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {

				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
					{
						URN:       "alert-history-odpf",
						Namespace: 2,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
					{
						URN:       "odpf-data-warning",
						Namespace: 1,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
								Configuration: map[string]string{
									"channel_name": "odpf-data",
								},
							},
						},
						Match: map[string]string{
							"environment": "integration",
							"team":        "odpf-data",
						},
					},
				}, nil)
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{
						ID: 1,
					},
				}, nil)
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
			},
			Provider: &provider.Provider{
				Type: "unknown",
			},
			ErrString: "subscriptions for provider type 'unknown' not supported",
		},
		{
			Description: "should return error if cortex client return error",
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
					{
						URN:       "alert-history-odpf",
						Namespace: 2,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
					{
						URN:       "odpf-data-warning",
						Namespace: 1,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
								Configuration: map[string]string{
									"channel_name": "odpf-data",
								},
							},
						},
						Match: map[string]string{
							"environment": "integration",
							"team":        "odpf-data",
						},
					},
				}, nil)
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{
						ID: 1,
					},
				}, nil)
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
				cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), mock.AnythingOfType("string")).Return(errors.New("some error"))
			},
			Provider: &provider.Provider{
				Type: provider.TypeCortex,
			},
			Namespace: &namespace.Namespace{
				ID:  2,
				URN: "namespace-urn",
			},
			ErrString: "error calling cortex: some error",
		},
		{
			Description: "should return nil error if cortex client return nil error",
			Setup: func(sr *mocks.SubscriptionRepository, ns *mocks.NamespaceService, ps *mocks.ProviderService, rs *mocks.ReceiverService, cc *mocks.CortexClient) {
				sr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("subscription.Filter")).Return([]subscription.Subscription{
					{
						URN:       "alert-history-odpf",
						Namespace: 2,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
					{
						URN:       "odpf-data-warning",
						Namespace: 1,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
								Configuration: map[string]string{
									"channel_name": "odpf-data",
								},
							},
						},
						Match: map[string]string{
							"environment": "integration",
							"team":        "odpf-data",
						},
					},
				}, nil)
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{
						ID: 1,
					},
				}, nil)
				rs.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{}, nil)
				cc.EXPECT().CreateAlertmanagerConfig(mock.AnythingOfType("cortex.AlertManagerConfig"), mock.AnythingOfType("string")).Return(nil)
			},
			Provider: &provider.Provider{
				Type: provider.TypeCortex,
			},
			Namespace: &namespace.Namespace{
				ID:  2,
				URN: "namespace-urn",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock       = new(mocks.SubscriptionRepository)
				namespaceServiceMock = new(mocks.NamespaceService)
				providerServiceMock  = new(mocks.ProviderService)
				receiverServiceMock  = new(mocks.ReceiverService)
				cortexClientMock     = new(mocks.CortexClient)
			)

			svc := subscription.NewService(repositoryMock, providerServiceMock, namespaceServiceMock, receiverServiceMock, cortexClientMock)

			tc.Setup(repositoryMock, namespaceServiceMock, providerServiceMock, receiverServiceMock, cortexClientMock)

			err := svc.SyncToUpstream(context.TODO(), tc.Namespace, tc.Provider)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			namespaceServiceMock.AssertExpectations(t)
			providerServiceMock.AssertExpectations(t)
			receiverServiceMock.AssertExpectations(t)
			cortexClientMock.AssertExpectations(t)
		})
	}
}
