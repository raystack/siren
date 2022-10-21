package slack_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/retry"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/plugins/receivers/slack/mocks"
	"github.com/stretchr/testify/mock"
)

func TestNotificationService_Publish(t *testing.T) {
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

			s := slack.NewNotificationService(mockSlackClient, nil)

			got, err := s.Publish(context.Background(), tt.notificationMessage)
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

func TestNotificationService_PreHookTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]interface{}
		want                  map[string]interface{}
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]interface{}{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token encryption failed",
			notificationConfigMap: map[string]interface{}{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return encrypted slack token if succeed",

			notificationConfigMap: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockEncryptor = new(mocks.Encryptor)
			)

			if tt.setup != nil {
				tt.setup(mockEncryptor)
			}

			s := slack.NewNotificationService(nil, mockEncryptor)
			got, err := s.PreHookTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationService.PreHookTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotificationService.PreHookTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationService_PostHookTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]interface{}
		want                  map[string]interface{}
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]interface{}{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token decryption failed",
			notificationConfigMap: map[string]interface{}{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return encrypted slack token if succeed",

			notificationConfigMap: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockEncryptor = new(mocks.Encryptor)
			)

			if tt.setup != nil {
				tt.setup(mockEncryptor)
			}

			s := slack.NewNotificationService(nil, mockEncryptor)
			got, err := s.PostHookTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationService.PostHookTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotificationService.PostHookTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackNotificationService_PreHookTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]interface{}
		want                  map[string]interface{}
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]interface{}{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token encryption failed",
			notificationConfigMap: map[string]interface{}{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return encrypted slack token if succeed",

			notificationConfigMap: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockEncryptor = new(mocks.Encryptor)
			)

			if tt.setup != nil {
				tt.setup(mockEncryptor)
			}

			s := slack.NewNotificationService(nil, mockEncryptor)
			got, err := s.PreHookTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackNotificationService.PreHookTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlackNotificationService.PreHookTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackNotificationService_PostHookTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]interface{}
		want                  map[string]interface{}
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]interface{}{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token decryption failed",
			notificationConfigMap: map[string]interface{}{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return encrypted slack token if succeed",

			notificationConfigMap: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				mockEncryptor = new(mocks.Encryptor)
			)

			if tt.setup != nil {
				tt.setup(mockEncryptor)
			}

			s := slack.NewNotificationService(nil, mockEncryptor)
			got, err := s.PostHookTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("SlackNotificationService.PostHookTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlackNotificationService.PostHookTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}
