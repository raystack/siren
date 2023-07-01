package notification_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/raystack/siren/core/log"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/core/notification/mocks"
	"github.com/raystack/siren/core/receiver"
	"github.com/raystack/siren/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func TestDispatchReceiverService_PrepareMessage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mocks.ReceiverService, *mocks.Notifier)
		n       notification.Notification
		want    []notification.Message
		want1   []log.Notification
		want2   bool
		wantErr bool
	}{
		{
			name: "should return error if receiver id in label is not parsable",
			n: notification.Notification{
				Labels: map[string]string{
					notification.ReceiverIDLabelKey: "x",
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if receiver service return error",
			n: notification.Notification{
				Labels: map[string]string{
					notification.ReceiverIDLabelKey: "11",
				},
			},
			setup: func(rs *mocks.ReceiverService, n *mocks.Notifier) {
				rs.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("receiver.GetOption")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return error if receiver type is unknown",
			n: notification.Notification{
				Labels: map[string]string{
					notification.ReceiverIDLabelKey: "11",
				},
			},
			setup: func(rs *mocks.ReceiverService, n *mocks.Notifier) {
				rs.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("receiver.GetOption")).Return(&receiver.Receiver{}, nil)
			},
			wantErr: true,
		},
		{
			name: "should return error if init message return error",
			n: notification.Notification{
				Labels: map[string]string{
					notification.ReceiverIDLabelKey: "11",
				},
			},
			setup: func(rs *mocks.ReceiverService, n *mocks.Notifier) {
				rs.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("receiver.GetOption")).Return(&receiver.Receiver{
					Type: testPluginType,
				}, nil)
				n.EXPECT().PreHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return no error if all flow passed",
			n: notification.Notification{
				Labels: map[string]string{
					notification.ReceiverIDLabelKey: "11",
				},
			},
			setup: func(rs *mocks.ReceiverService, n *mocks.Notifier) {
				rs.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64"), mock.AnythingOfType("receiver.GetOption")).Return(&receiver.Receiver{
					ID:   11,
					Type: testPluginType,
				}, nil)
				n.EXPECT().PreHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]interface{}{}, nil)
			},
			want: []notification.Message{
				{
					Status:       notification.MessageStatusEnqueued,
					ReceiverType: testPluginType,
					Configs:      map[string]interface{}{},
					Details:      map[string]interface{}{"notification_type": string(""), "receiver_id": string("11")},
					MaxTries:     3,
				},
			},
			want1: []log.Notification{{ReceiverID: 11}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockReceiverService = new(mocks.ReceiverService)
				mockNotifier        = new(mocks.Notifier)
			)
			s := notification.NewDispatchReceiverService(
				mockReceiverService,
				map[string]notification.Notifier{
					testPluginType: mockNotifier,
				})
			if tt.setup != nil {
				tt.setup(mockReceiverService, mockNotifier)
			}
			got, got1, got2, err := s.PrepareMessage(context.TODO(), tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("DispatchReceiverService.PrepareMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want,
				cmpopts.IgnoreFields(notification.Message{}, "ID", "CreatedAt", "UpdatedAt"),
				cmpopts.IgnoreUnexported(notification.Message{})); diff != "" {
				t.Errorf("DispatchReceiverService.PrepareMessage() diff = %v", diff)
			}
			if diff := cmp.Diff(got1, tt.want1); diff != "" {
				t.Errorf("DispatchReceiverService.PrepareMessage() diff = %v", diff)
			}
			if got2 != tt.want2 {
				t.Errorf("DispatchReceiverService.PrepareMessage() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
