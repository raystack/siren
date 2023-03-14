package slack_test

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
	"github.com/stretchr/testify/mock"
)

func TestService_BuildData(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(sc *mocks.SlackCaller, e *mocks.Encryptor)
		Confs       map[string]interface{}
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if configuration is invalid",
				Setup:       func(sc *mocks.SlackCaller, e *mocks.Encryptor) {},
				Confs:       make(map[string]interface{}),
				Err:         errors.New("invalid slack receiver config, workspace: , token: "),
			},
			{
				Description: "should return error if failed to get workspace channels with slack client",
				Setup: func(sc *mocks.SlackCaller, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("secret.MaskableString")).Return(nil, errors.New("some error"))
				},
				Confs: map[string]interface{}{
					"token":     secret.MaskableString("key"),
					"workspace": "gotocompany",
				},
				Err: errors.New("could not get channels: some error"),
			},
			{
				Description: "should return nil error if success populating receiver.Receiver",
				Setup: func(sc *mocks.SlackCaller, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("secret.MaskableString")).Return([]slack.Channel{
						{
							ID:   "id",
							Name: "name",
						},
					}, nil)
				},
				Confs: map[string]interface{}{
					"token":     secret.MaskableString("key"),
					"workspace": "gotocompany",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				slackClientMock = new(mocks.SlackCaller)
				encryptorMock   = new(mocks.Encryptor)
			)

			svc := slack.NewPluginService(slack.AppConfig{}, encryptorMock, slack.WithSlackClient(slackClientMock))

			tc.Setup(slackClientMock, encryptorMock)

			_, err := svc.BuildData(ctx, tc.Confs)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			slackClientMock.AssertExpectations(t)
			encryptorMock.AssertExpectations(t)
		})
	}
}

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

			s := slack.NewPluginService(slack.AppConfig{}, nil, slack.WithSlackClient(mockSlackClient))

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
				"channel_type": "",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Encrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
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

			s := slack.NewPluginService(slack.AppConfig{}, mockEncryptor)
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
			name: "should return decrypted slack token if succeed",

			notificationConfigMap: map[string]interface{}{
				"workspace":    "a workspace",
				"token":        secret.MaskableString("a token"),
				"channel_name": "channel",
				"channel_type": "",
			},
			setup: func(e *mocks.Encryptor) {
				e.EXPECT().Decrypt(mock.AnythingOfType("secret.MaskableString")).Return(secret.MaskableString("maskable-token"), nil)
			},
			want: map[string]interface{}{
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

			s := slack.NewPluginService(slack.AppConfig{}, mockEncryptor)
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
