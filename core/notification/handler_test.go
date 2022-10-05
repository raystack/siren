package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/notification/mocks"
	"github.com/stretchr/testify/mock"
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
			wantErr: true,
		},
		{
			name: "return error if publish message return error and error handler queue return error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				n.EXPECT().Publish(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, errors.New("some error"))
				q.EXPECT().ErrorHandler(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "return error if publish message return error and error handler queue return no error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				n.EXPECT().Publish(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, errors.New("some error"))
				q.EXPECT().ErrorHandler(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "return error if publish message success and success handler queue return error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				n.EXPECT().Publish(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, nil)
				q.EXPECT().SuccessHandler(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "return no error if publish message success and success handler queue return no error",
			messages: []notification.Message{
				{
					ReceiverType: testPluginType,
				},
			},
			setup: func(q *mocks.Queuer, n *mocks.Notifier) {
				n.EXPECT().Publish(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(false, nil)
				q.EXPECT().SuccessHandler(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("notification.Message")).Return(nil)
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

			h := notification.NewHandler(log.NewNoop(), mockQueue, map[string]notification.Notifier{
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
