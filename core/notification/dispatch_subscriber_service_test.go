package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	saltlog "github.com/goto/salt/log"
	"github.com/goto/siren/core/log"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/notification/mocks"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/stretchr/testify/mock"
)

func TestDispatchSubscriberService_PrepareMessage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mocks.SubscriptionService, *mocks.SilenceService, *mocks.Notifier)
		n       notification.Notification
		want    []notification.Message
		want1   []log.Notification
		want2   bool
		wantErr bool
	}{
		{
			name: "should return error if subscription service match by labels return error",
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return error if no matching subscriptions",
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "should return error if match subscription exist but list silences return error",
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID: 123,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("silence.Filter")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return error if match subscription exist but list silences by label return error",
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID: 123,
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("silence.Filter")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return no error if silenced by labels success",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID: 123,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return([]silence.Silence{
					{
						ID:          "silence-id",
						NamespaceID: 1,
						TargetID:    123,
					},
				}, nil)
			},
			want:  []notification.Message{},
			want1: []log.Notification{{SubscriptionID: 123, NamespaceID: 1, SilenceIDs: []string{"silence-id"}}},
			want2: true,
		},
		{
			name: "should return error if silenced by subscription return error",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID:        123,
						Namespace: 1,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return(nil, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID:    1,
					SubscriptionID: 123,
				}).Return([]silence.Silence{
					{
						ID:          "silence-id",
						NamespaceID: 1,
					},
				}, nil)
			},
			wantErr: true,
		},
		{
			name: "should return no error if silenced by subscription success",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID:        123,
						Namespace: 1,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return(nil, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID:    1,
					SubscriptionID: 123,
				}).Return([]silence.Silence{
					{
						ID:          "silence-id",
						NamespaceID: 1,
						Type:        silence.TypeSubscription,
					},
				}, nil)
			},
			want:  []notification.Message{},
			want1: []log.Notification{{SubscriptionID: 123, NamespaceID: 1, ReceiverID: 1, SilenceIDs: []string{"silence-id"}}},
			want2: true,
		},
		{
			name: "should return error if receiver type is unknown",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID:        123,
						Namespace: 1,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID: 1,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return(nil, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID:    1,
					SubscriptionID: 123,
				}).Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name: "should return error if init messages return error",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID:        123,
						Namespace: 1,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID:   1,
								Type: testPluginType,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return(nil, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID:    1,
					SubscriptionID: 123,
				}).Return(nil, nil)
				n.EXPECT().PreHookQueueTransformConfigs(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return no error if all flow passed and no silences",
			n: notification.Notification{
				NamespaceID: 1,
			},
			setup: func(ss1 *mocks.SubscriptionService, ss2 *mocks.SilenceService, n *mocks.Notifier) {
				ss1.EXPECT().MatchByLabels(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("map[string]string")).Return([]subscription.Subscription{
					{
						ID:        123,
						Namespace: 1,
						Match: map[string]string{
							"k1": "v1",
						},
						Receivers: []subscription.Receiver{
							{
								ID:   1,
								Type: testPluginType,
							},
						},
					},
				}, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID: 1,
					SubscriptionMatch: map[string]string{
						"k1": "v1",
					},
				}).Return(nil, nil)
				ss2.EXPECT().List(mock.AnythingOfType("context.todoCtx"), silence.Filter{
					NamespaceID:    1,
					SubscriptionID: 123,
				}).Return(nil, nil)
				n.EXPECT().PreHookQueueTransformConfigs(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]any{}, nil)
			},
			want: []notification.Message{
				{
					Status:       notification.MessageStatusEnqueued,
					ReceiverType: testPluginType,
					Configs:      map[string]any{},
					Details:      map[string]any{"notification_type": string("")},
					MaxTries:     3,
				},
			},
			want1: []log.Notification{{NamespaceID: 1, SubscriptionID: 123, ReceiverID: 1}},
			want2: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockSubscriptionService = new(mocks.SubscriptionService)
				mockSilenceService      = new(mocks.SilenceService)
				mockNotifier            = new(mocks.Notifier)
			)
			s := notification.NewDispatchSubscriberService(
				saltlog.NewNoop(),
				mockSubscriptionService,
				mockSilenceService, map[string]notification.Notifier{
					testPluginType: mockNotifier,
				},
				true,
			)

			if tt.setup != nil {
				tt.setup(mockSubscriptionService, mockSilenceService, mockNotifier)
			}

			got, got1, got2, err := s.PrepareMessage(context.TODO(), tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("DispatchSubscriberService.PrepareMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want,
				cmpopts.IgnoreFields(notification.Message{}, "ID", "CreatedAt", "UpdatedAt"),
				cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
				t.Errorf("DispatchSubscriberService.PrepareMessage() diff = %v", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("DispatchSubscriberService.PrepareMessage() diff = %v", diff)
			}
			if got2 != tt.want2 {
				t.Errorf("DispatchSubscriberService.PrepareMessage() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
