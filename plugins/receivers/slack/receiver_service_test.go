package slack_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/plugins/receivers/slack/mocks"
	"github.com/stretchr/testify/mock"
)

func TestReceiverService_BuildData(t *testing.T) {
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
					"workspace": "odpf",
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
					"workspace": "odpf",
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

			svc := slack.NewReceiverService(slackClientMock, encryptorMock)

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

func TestReceiverService_BuildNotificationConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]interface{}
		ReceiverConfigs     map[string]interface{}
		ExpectedConfigMap   map[string]interface{}
		wantErr             bool
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if receiver 'token' exist but it is not string",
				SubscriptionConfigs: map[string]interface{}{
					"channel_name": "odpf_warning",
				},
				ReceiverConfigs: map[string]interface{}{
					"token": 123,
				},
				wantErr: true,
			},
			{
				Description: "should return configs without token if receiver 'token' does not exist",
				SubscriptionConfigs: map[string]interface{}{
					"channel_name": "odpf_warning",
				},
				ExpectedConfigMap: map[string]interface{}{
					"channel_name": "odpf_warning",
					"token":        secret.MaskableString(""),
					"workspace":    "",
				},
			},
			{
				Description: "should return configs with token if receiver 'token' exist in string",
				SubscriptionConfigs: map[string]interface{}{
					"channel_name": "odpf_warning",
				},
				ReceiverConfigs: map[string]interface{}{
					"token": "123",
				},
				ExpectedConfigMap: map[string]interface{}{
					"channel_name": "odpf_warning",
					"token":        secret.MaskableString("123"),
					"workspace":    "",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := slack.NewReceiverService(nil, nil)

			got, err := svc.BuildNotificationConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
			if (err != nil) != tc.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil {
				if !cmp.Equal(got, tc.ExpectedConfigMap) {
					t.Errorf("got result %+v, expected was %+v", got, tc.ExpectedConfigMap)
				}
			}
		})
	}
}
