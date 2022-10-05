package httpreceiver_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/plugins/receivers/httpreceiver"
	"github.com/odpf/siren/plugins/receivers/httpreceiver/mocks"
	"github.com/stretchr/testify/mock"
)

func TestHTTPNotificationService_Publish(t *testing.T) {
	tests := []struct {
		name                string
		setup               func(*mocks.HTTPCaller)
		notificationMessage notification.Message
		wantRetryable       bool
		wantErr             bool
	}{
		{
			name: "should return error if failed to decode notification config",
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"url": true,
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if failed to decode notification detail",
			notificationMessage: notification.Message{
				Detail: map[string]interface{}{
					"description": make(chan bool),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error and not retryable if notify return error",
			setup: func(pd *mocks.HTTPCaller) {
				pd.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(errors.New("some error"))
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"url": "123123",
				},
				Detail: map[string]interface{}{
					"description": "hello",
				},
			},
			wantRetryable: false,
			wantErr:       true,
		},
		{
			name: "should return error and retryable if notify return retryable error",
			setup: func(sc *mocks.HTTPCaller) {
				sc.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(retry.RetryableError{Err: errors.New("some error")})
			},
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"url": "123123",
				},
				Detail: map[string]interface{}{
					"description": "hello",
				},
			},
			wantRetryable: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPCaller := new(mocks.HTTPCaller)

			if tt.setup != nil {
				tt.setup(mockHTTPCaller)
			}

			pd := httpreceiver.NewNotificationService(mockHTTPCaller)

			got, err := pd.Publish(context.Background(), tt.notificationMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPNotificationService.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRetryable {
				t.Errorf("HTTPNotificationService.Publish() = %v, want %v", got, tt.wantRetryable)
			}
		})
	}
}
