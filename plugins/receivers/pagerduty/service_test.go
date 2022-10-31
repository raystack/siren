package pagerduty_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/plugins/receivers/pagerduty"
	"github.com/odpf/siren/plugins/receivers/pagerduty/mocks"
	"github.com/stretchr/testify/mock"
)

func TestService_Send_V1(t *testing.T) {
	tests := []struct {
		name                string
		setup               func(*mocks.PagerDutyCaller)
		notificationMessage notification.Message
		wantRetryable       bool
		wantErr             bool
	}{
		{
			name: "should return error if failed to decode notification config",
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"service_key": true,
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if failed to decode notification detail",
			notificationMessage: notification.Message{
				Details: map[string]interface{}{
					"description": make(chan bool),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error and not retryable if notify return error",
			setup: func(pd *mocks.PagerDutyCaller) {
				pd.EXPECT().NotifyV1(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("pagerduty.MessageV1")).Return(errors.New("some error"))
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"service_key": "123123",
				},
				Details: map[string]interface{}{
					"description": "hello",
				},
			},
			wantRetryable: false,
			wantErr:       true,
		},
		{
			name: "should return error and retryable if notify return retryable error",
			setup: func(sc *mocks.PagerDutyCaller) {
				sc.EXPECT().NotifyV1(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("pagerduty.MessageV1")).Return(retry.RetryableError{Err: errors.New("some error")})
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"service_key": "123123",
				},
				Details: map[string]interface{}{
					"description": "hello",
				},
			},
			wantRetryable: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPDClient := new(mocks.PagerDutyCaller)

			if tt.setup != nil {
				tt.setup(mockPDClient)
			}

			pd := pagerduty.NewPluginService(pagerduty.AppConfig{}, pagerduty.WithPagerDutyClient(mockPDClient))

			got, err := pd.Send(context.Background(), tt.notificationMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationService.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRetryable {
				t.Errorf("NotificationService.Send() = %v, want %v", got, tt.wantRetryable)
			}
		})
	}
}
