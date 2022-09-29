package file_test

import (
	"context"
	"testing"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/plugins/receivers/file"
)

func TestNotificationService_Publish(t *testing.T) {
	tests := []struct {
		name                string
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
				Details: map[string]interface{}{
					"description": make(chan bool),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error and not retryable if notify return error",
			notificationMessage: notification.Message{
				Configs: map[string]interface{}{
					"url": "123123",
				},
				Details: map[string]interface{}{
					"description": "hello",
				},
			},
			wantRetryable: false,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fr := file.NewNotificationService()

			got, err := fr.Publish(context.Background(), tt.notificationMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationService.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRetryable {
				t.Errorf("NotificationService.Publish() = %v, want %v", got, tt.wantRetryable)
			}
		})
	}
}
