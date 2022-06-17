package receiver_test

import (
	"errors"
	"testing"
	"time"

	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
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
				Err:         errors.New("no token field found"),
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
				Err:         errors.New("no token field found"),
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
		testCases = []testCase{
			{
				Description: "should return error if no token field in configurations",
				Setup:       func(sc *mocks.SlackClient, e *mocks.Encryptor) {},
				Rcv:         &receiver.Receiver{},
				Err:         errors.New("no token found in configurations"),
			},
			{
				Description: "should return error if failed to get workspace channels with slack client",
				Setup: func(sc *mocks.SlackClient, e *mocks.Encryptor) {
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("slack.ClientCallOption"), mock.AnythingOfType("slack.ClientCallOption")).Return(nil, errors.New("some error"))
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
					sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("slack.ClientCallOption"), mock.AnythingOfType("slack.ClientCallOption")).Return([]slack.Channel{
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

			_, err := svc.PopulateReceiver(tc.Rcv)
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
