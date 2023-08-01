package slackchannel_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/pkg/secret"
	"github.com/goto/siren/plugins/receivers/slack"
	"github.com/goto/siren/plugins/receivers/slack/mocks"
	"github.com/goto/siren/plugins/receivers/slackchannel"
	"github.com/stretchr/testify/mock"
)

func TestService_Send(t *testing.T) {
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
				Configs: map[string]any{
					"token": true,
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if failed to decode notification detail",
			notificationMessage: notification.Message{
				Details: map[string]any{
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
				Configs: map[string]any{
					"token": true,
				},
				Details: map[string]any{
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
				Configs: map[string]any{
					"token": "123123",
				},
				Details: map[string]any{
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

			s := slackchannel.NewPluginService(slack.AppConfig{}, nil, slack.WithSlackClient(mockSlackClient))

			got, err := s.Send(context.Background(), tt.notificationMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRetryable {
				t.Errorf("Service.Publish() = %v, want %v", got, tt.wantRetryable)
			}
		})
	}
}

func TestService_PreHookQueueTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]any
		want                  map[string]any
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]any{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token encryption failed",
			notificationConfigMap: map[string]any{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return encrypted slack token if succeed",

			notificationConfigMap: map[string]any{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
				"channel_type": "",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]any{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
				"channel_type": "",
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

			s := slackchannel.NewPluginService(slack.AppConfig{}, mockEncryptor)
			got, err := s.PreHookQueueTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PreHookQueueTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.PreHookQueueTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_PostHookQueueTransformConfigs(t *testing.T) {
	tests := []struct {
		name                  string
		setup                 func(*mocks.Encryptor)
		notificationConfigMap map[string]any
		want                  map[string]any
		wantErr               bool
	}{
		{
			name:                  "should return error if failed to parse configmap to notification config",
			notificationConfigMap: nil,
			wantErr:               true,
		},
		{
			name: "should return error if validate notification config failed",
			notificationConfigMap: map[string]any{
				"token": 123,
			},
			wantErr: true,
		},
		{
			name: "should return error if slack token decryption failed",
			notificationConfigMap: map[string]any{
				"token": secret.MaskableString("a token"),
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return("", errors.New("some error"))
			},
			wantErr: true,
		},
		{
			name: "should return decrypted slack token if succeed",

			notificationConfigMap: map[string]any{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
				"channel_type": "",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]any{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("maskable-token"),
				"channel_name": "channel",
				"channel_type": "",
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

			s := slackchannel.NewPluginService(slack.AppConfig{}, mockEncryptor)
			got, err := s.PostHookQueueTransformConfigs(context.TODO(), tt.notificationConfigMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PostHookQueueTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.PostHookQueueTransformConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluginService_PreHookDBTransformConfigs(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*mocks.Encryptor)
		configurations map[string]any
		wantErr        bool
	}{
		{
			name:    "should return error if channel_name is missing",
			wantErr: true,
		},
		{
			name: "shouldd return non error if channel_name is not missing",
			configurations: map[string]any{
				"channel_name": "a-channel",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &slackchannel.PluginService{}
			_, err := s.PreHookDBTransformConfigs(context.TODO(), tt.configurations)
			if (err != nil) != tt.wantErr {
				t.Errorf("PluginService.PreHookDBTransformConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
