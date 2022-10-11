package slack_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/plugins/receivers/slack/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSlackNotificationService_Publish(t *testing.T) {
	tests := []struct {
		name                string
		setup               func(*mocks.SlackCaller)
		notificationMessage notification.Message
		wantRetryable       bool
		wantErr             bool
	}{
		{
			name: "should return error if failed to decode notification config",
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"token": true,
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if failed to decode notification detail",
			notificationMessage: notification.Message{
				Details: map[string]interface{}{
					"text": make(chan bool),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error and not retryable if notify return error",
			setup: func(sc *mocks.SlackCaller) {
				sc.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("slack.NotificationConfig"), mock.AnythingOfType("slack.Message")).Return(errors.New("some error"))
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"token": true,
				},
				Details: map[string]interface{}{
					"message": "hello",
				},
			},
			wantRetryable: false,
			wantErr:       true,
		},
		{
			name: "should return error and retryable if notify return retryable error",
			setup: func(sc *mocks.SlackCaller) {
				sc.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("slack.NotificationConfig"), mock.AnythingOfType("slack.Message")).Return(retry.RetryableError{Err: errors.New("some error")})
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"token": "123123",
				},
				Details: map[string]interface{}{
					"message": "hello",
				},
			},
			wantRetryable: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSlackClient := new(mocks.SlackCaller)

			if tt.setup != nil {
				tt.setup(mockSlackClient)
			}

			s := slack.NewNotificationService(mockSlackClient)

			got, err := s.Publish(context.Background(), tt.notificationMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackNotificationService.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRetryable {
				t.Errorf("SlackNotificationService.Publish() = %v, want %v", got, tt.wantRetryable)
			}
		})
	}
}
