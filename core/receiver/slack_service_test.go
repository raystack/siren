package receiver_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/slack"
	"github.com/stretchr/testify/mock"
)

func TestSlackService_Encrypt(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		Setup       func(*mocks.SlackClient, *mocks.Encryptor)
		Err         error
	}

	var (
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return error if no token in configurations field in encrypt error",
				Setup:       func(sc *mocks.SlackClient, e *mocks.Encryptor) {},
				Rcv:         &receiver.Receiver{},
				Err:         errors.New("no token in configurations found"),
			},
			{
				Description: "should return error if encrypt error",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", errors.New("encrypt error"))
				},
				Rcv: &receiver.Receiver{
					Configurations: map[string]interface{}{
						"token": "key",
					}},
				Err: errors.New("encrypt error"),
			},
			{
				Description: "should success if encrypt success",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				},
				Rcv: &receiver.Receiver{
					ID:   10,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Configurations: map[string]interface{}{
						"token": "key",
					},
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				slackClientMock = new(mocks.SlackClient)
				encryptorMock   = new(mocks.Encryptor)
			)

			svc := receiver.NewSlackService(slackClientMock, encryptorMock)

			tc.Setup(slackClientMock, encryptorMock)

			err := svc.Encrypt(tc.Rcv)
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

func TestSlackService_Decrypt(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		Setup       func(*mocks.SlackClient, *mocks.Encryptor)
		Err         error
	}

	var (
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return error if no token in configurations field in decrypt error",
				Setup:       func(sc *mocks.SlackClient, e *mocks.Encryptor) {},
				Rcv:         &receiver.Receiver{},
				Err:         errors.New("no token in configurations found"),
			},
			{
				Description: "should return error if decrypt error",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", errors.New("decrypt error"))
				},
				Rcv: &receiver.Receiver{
					Configurations: map[string]interface{}{
						"token": "key",
					}},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if decrypt success",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", nil)
				},
				Rcv: &receiver.Receiver{
					ID:   10,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Configurations: map[string]interface{}{
						"token": "key",
					},
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				slackClientMock = new(mocks.SlackClient)
				encryptorMock   = new(mocks.Encryptor)
			)

			svc := receiver.NewSlackService(slackClientMock, encryptorMock)

			tc.Setup(slackClientMock, encryptorMock)

			err := svc.Decrypt(tc.Rcv)
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

func TestSlackService_PopulateReceiver(t *testing.T) {

	type testCase struct {
		Description string
		Setup       func(sc *mocks.SlackClient, e *mocks.Encryptor)
		Rcv         *receiver.Receiver
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if no token field in configurations",
				Setup:       func(sc *mocks.SlackClient, e *mocks.Encryptor) {},
				Rcv:         &receiver.Receiver{},
				Err:         errors.New("no token in configurations found"),
			},
			{
				Description: "should return error if failed to get workspace channels with slack client",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("slack.ClientCallOption")).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					Configurations: map[string]interface{}{
						"token": "key",
					}},
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
				Rcv: &receiver.Receiver{
					Configurations: map[string]interface{}{
						"token": "key",
					}},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				slackClientMock = new(mocks.SlackClient)
				encryptorMock   = new(mocks.Encryptor)
			)

			svc := receiver.NewSlackService(slackClientMock, encryptorMock)

			tc.Setup(slackClientMock, encryptorMock)

			_, err := svc.PopulateReceiver(ctx, tc.Rcv)
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

func TestSlackService_ValidateConfiguration(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		ErrString   string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if 'client_id' is empty",
				Rcv:         &receiver.Receiver{},
				ErrString:   "no value supplied for required configurations map key \"client_id\"",
			},
			{
				Description: "should return error if 'client_secret' is empty",
				Rcv: &receiver.Receiver{
					Configurations: receiver.Configurations{
						"client_id": "client_id",
					},
				},
				ErrString: "no value supplied for required configurations map key \"client_secret\"",
			},
			{
				Description: "should return error if 'auth_code' is empty",
				Rcv: &receiver.Receiver{
					Configurations: receiver.Configurations{
						"client_id":     "client_id",
						"client_secret": "client_secret",
					},
				},
				ErrString: "no value supplied for required configurations map key \"auth_code\"",
			},
			{
				Description: "should return nil error if all configurations are valid",
				Rcv: &receiver.Receiver{
					Configurations: receiver.Configurations{
						"client_id":     "client_id",
						"client_secret": "client_secret",
						"auth_code":     "auth_code",
					},
				},
			},
			{
				Description: "should return error if receiver is nil",
				ErrString:   "receiver to validate is nil",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := receiver.NewSlackService(nil, nil)

			err := svc.ValidateConfiguration(tc.Rcv)
			if err != nil {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestSlackService_GetSubscriptionConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]string
		ReceiverConfigs     receiver.Configurations
		ExpectedConfigMap   map[string]string
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
				SubscriptionConfigs: map[string]string{
					"channel_name": "odpf_warning",
				},
				ReceiverConfigs: receiver.Configurations{
					"token": 123,
				},
				ErrString: "token config from receiver should be in string",
			},
			{
				Description: "should return configs without token if receiver 'token' does not exist", //TODO might need to check this behaviour, should be returning error
				SubscriptionConfigs: map[string]string{
					"channel_name": "odpf_warning",
				},
				ExpectedConfigMap: map[string]string{
					"channel_name": "odpf_warning",
				},
			},
			{
				Description: "should return configs with token if receiver 'token'exist in string",
				SubscriptionConfigs: map[string]string{
					"channel_name": "odpf_warning",
				},
				ReceiverConfigs: receiver.Configurations{
					"token": "123",
				},
				ExpectedConfigMap: map[string]string{
					"channel_name": "odpf_warning",
					"token":        "123",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := receiver.NewSlackService(nil, nil)

			got, err := svc.GetSubscriptionConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
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
