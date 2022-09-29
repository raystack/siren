package slack_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/plugins/receivers/slack/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSlackReceiverService_BuildData(t *testing.T) {

	type testCase struct {
		Description string
		Setup       func(sc *mocks.SlackClient, e *mocks.Encryptor)
		Confs       map[string]interface{}
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if no token field in configurations",
				Setup:       func(sc *mocks.SlackClient, e *mocks.Encryptor) {},
				Confs:       make(map[string]interface{}),
				Err:         errors.New("no token in configurations found"),
			},
			{
				Description: "should return error if failed to get workspace channels with slack client",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("slack.ClientCallOption")).Return(nil, errors.New("some error"))
				},
				Confs: map[string]interface{}{
					"token": "key",
				},
				Err: errors.New("could not get channels: some error"),
			},
			{
				Description: "should return nil error if success populating receiver.Receiver",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("slack.ClientCallOption")).Return([]slack.Channel{
						{
							ID:   "id",
							Name: "name",
						},
					}, nil)
				},
				Confs: map[string]interface{}{
					"token": "key",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				slackClientMock = new(mocks.SlackClient)
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

func TestSlackReceiverService_Notify(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.SlackClient)
		Confs       map[string]interface{}
		Message     map[string]interface{}
		ErrString   string
	}
	var (
		testCases = []testCase{
			{
				Description: "should return error if cannot cast token to string",
				Setup:       func(sc *mocks.SlackClient) {},
				Confs:       make(map[string]interface{}),
				ErrString:   "no token in configurations found",
			},
			{
				Description: "should return error if message cannot be converted to slack message",
				Setup:       func(sc *mocks.SlackClient) {},
				Confs: map[string]interface{}{
					"token": "123123",
				},
				ErrString: "non empty message or non zero length block is required",
			},
			{
				Description: "should return error if slack client return error",
				Setup: func(sc *mocks.SlackClient) {
					sc.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*slack.Message"), mock.AnythingOfType("slack.ClientCallOption")).Return(errors.New("some error"))
				},
				Confs: map[string]interface{}{
					"token": "123123",
				},
				Message: map[string]interface{}{
					"receiver_name": "receiver_name",
					"receiver_type": "channel",
					"message":       "message",
					"blocks": []map[string]interface{}{
						{
							"type": "section",
						},
					},
				},
				ErrString: "error calling slack notify: some error",
			},
			{
				Description: "should return nil error if slack client return nil error",
				Setup: func(sc *mocks.SlackClient) {
					sc.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*slack.Message"), mock.AnythingOfType("slack.ClientCallOption")).Return(nil)
				},
				Confs: map[string]interface{}{
					"token": "123123",
				},
				Message: map[string]interface{}{
					"receiver_name": "receiver_name",
					"receiver_type": "channel",
					"message":       "message",
					"blocks": []map[string]interface{}{
						{
							"type": "section",
						},
					},
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				clientMock = new(mocks.SlackClient)
			)

			svc := slack.NewReceiverService(clientMock, nil)

			tc.Setup(clientMock)

			err := svc.Notify(context.TODO(), tc.Confs, tc.Message)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			clientMock.AssertExpectations(t)
		})
	}
}

// func TestSlackReceiverService_ValidateConfigurations(t *testing.T) {
// 	type testCase struct {
// 		Description string
// 		Confs       map[string]interface{}
// 		ErrString   string
// 	}

// 	var (
// 		testCases = []testCase{
// 			{
// 				Description: "should return error if 'client_id' is empty",
// 				Confs:       make(map[string]interface{}),
// 				ErrString:   "no value supplied for required configurations map key \"client_id\"",
// 			},
// 			{
// 				Description: "should return error if 'client_secret' is empty",
// 				Confs: map[string]interface{}{
// 					"client_id": "client_id",
// 				},
// 				ErrString: "no value supplied for required configurations map key \"client_secret\"",
// 			},
// 			{
// 				Description: "should return error if 'auth_code' is empty",
// 				Confs: map[string]interface{}{
// 					"client_id":     "client_id",
// 					"client_secret": "client_secret",
// 				},
// 				ErrString: "no value supplied for required configurations map key \"auth_code\"",
// 			},
// 			{
// 				Description: "should return nil error if all configurations are valid",
// 				Confs: map[string]interface{}{
// 					"client_id":     "client_id",
// 					"client_secret": "client_secret",
// 					"auth_code":     "auth_code",
// 				},
// 			},
// 			{
// 				Description: "should return nil error if a configuration is not in string",
// 				Confs: map[string]interface{}{
// 					"client_id":     123,
// 					"client_secret": "client_secret",
// 					"auth_code":     "auth_code",
// 				},
// 				ErrString: "wrong type for configurations map key \"client_id\": expected type string, got value 123 of type int",
// 			},
// 		}
// 	)

// 	for _, tc := range testCases {
// 		t.Run(tc.Description, func(t *testing.T) {
// 			svc := slack.NewReceiverService(nil, nil)

// 			err := svc.ValidateConfigurations(tc.Confs)
// 			if err != nil {
// 				if tc.ErrString != err.Error() {
// 					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
// 				}
// 			}
// 		})
// 	}
// }

func TestSlackReceiverService_EnrichSubscriptionConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]interface{}
		ReceiverConfigs     map[string]interface{}
		ExpectedConfigMap   map[string]interface{}
		ErrString           string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if 'channel_name' does not exist",
				ErrString:   "subscription receiver config 'channel_name' was missing",
			},
			{
				Description: "should return error if receiver 'token' exist but it is not string",
				SubscriptionConfigs: map[string]interface{}{
					"channel_name": "odpf_warning",
				},
				ReceiverConfigs: map[string]interface{}{
					"token": 123,
				},
				ErrString: "token config from receiver should be in string",
			},
			{
				Description: "should return configs without token if receiver 'token' does not exist", //TODO might need to check this behaviour, should be returning error
				SubscriptionConfigs: map[string]interface{}{
					"channel_name": "odpf_warning",
				},
				ExpectedConfigMap: map[string]interface{}{
					"channel_name": "odpf_warning",
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
					"token":        "123",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := slack.NewReceiverService(nil, nil)

			got, err := svc.BuildNotificationConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedConfigMap) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedConfigMap)
			}
		})
	}
}
