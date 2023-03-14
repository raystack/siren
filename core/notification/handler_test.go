package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/salt/log"
	"github.com/stretchr/testify/mock"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/notification/mocks"
)

const testReceiverType = "test"

func TestHandler_MessageHandler(t *testing.T) {
	testCases := []struct {
		name     string
		messages []notification.Message
		setup    func(*mocks.Queuer, *mocks.Notifier)
		wantErr  bool
	}{
		{
			name: "return error if plugin type is not supported",
			messages: []notification.Message{
				{
					ReceiverType: "random",
				},
			},
			setup: func(q *mocks.Queuer, _ *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
			},
			wantErr: true,
		},
		{
			name: "return error if post hook transform config is failing and error callback success",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("some error"))
				q.EXPECT().ErrorCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "return error if post hook transform config is failing and error callback is failing",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("some error"))
				q.EXPECT().ErrorCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "return error if send message return error and error handler queue return error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]interface{}{}, nil)
				n.EXPECT().Send(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, errors.New("some error"))
				q.EXPECT().ErrorCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "return error if send message return error and error handler queue return no error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]interface{}{}, nil)
				n.EXPECT().Send(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, errors.New("some error"))
				q.EXPECT().ErrorCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "return error if send message success and success handler queue return error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]interface{}{}, nil)
				n.EXPECT().Send(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, nil)
				q.EXPECT().SuccessCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "return no error if send message success and success handler queue return no error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				q.EXPECT().Type().Return("postgresql")
				n.EXPECT().PostHookQueueTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("map[string]interface {}")).Return(map[string]interface{}{}, nil)
				n.EXPECT().Send(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, nil)
				q.EXPECT().SuccessCallback(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				mockQueue    = new(mocks.Queuer)
				mockNotifier = new(mocks.Notifier)
			)

			if tc.setup != nil {
				tc.setup(mockQueue, mockNotifier)
			}

			h := notification.NewHandler(notification.HandlerConfig{}, log.NewNoop(), mockQueue, map[string]notification.Notifier{
				testReceiverType: mockNotifier,
			})
			if err := h.MessageHandler(context.TODO(), tc.messages); (err != nil) != tc.wantErr {
				t.Errorf("Handler.messageHandler() error = %v, wantErr %v", err, tc.wantErr)
			}

			mockQueue.AssertExpectations(t)
			mockNotifier.AssertExpectations(t)
		})
	}
}
