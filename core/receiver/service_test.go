package receiver_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	"github.com/odpf/siren/pkg/slack"
	"github.com/stretchr/testify/mock"
)

func TestService_ListReceivers_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Receivers   []*receiver.Receiver
		Setup       func(*mocks.SecureServiceProxy, testCase)
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should call service List method and return result in domain's type",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().ListReceivers().Return(tc.Receivers, nil)
			},
			Receivers: []*receiver.Receiver{
				{
					Id:   10,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Configurations: map[string]interface{}{
						"token": "key",
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Err: nil,
		},
		{
			Description: "should call service List method and return error if any",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().ListReceivers().Return(tc.Receivers, tc.Err)
			},
			Err: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				secureServiceMock = new(mocks.SecureServiceProxy)
				slackClientMock   = new(mocks.SlackClient)
			)

			svc := receiver.NewService(secureServiceMock, slackClientMock)

			tc.Setup(secureServiceMock, tc)

			got, err := svc.ListReceivers()
			if tc.Err != err {
				t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
			}
			if !cmp.Equal(got, tc.Receivers) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.Receivers)
			}
			secureServiceMock.AssertExpectations(t)
			slackClientMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateReceiver_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.SecureServiceProxy, testCase)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should call repository Create method and return nil error",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().CreateReceiver(tc.Rcv).Return(nil)
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
			},
			Err: nil,
		},
		{
			Description: "should call repository Create method and return error",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().CreateReceiver(tc.Rcv).Return(tc.Err)
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
			},
			Err: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				secureServiceMock = new(mocks.SecureServiceProxy)
				slackClientMock   = new(mocks.SlackClient)
			)
			svc := receiver.NewService(secureServiceMock, slackClientMock)
			tc.Setup(secureServiceMock, tc)

			err := svc.CreateReceiver(tc.Rcv)
			if tc.Err != err {
				t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
			}

			secureServiceMock.AssertExpectations(t)
			slackClientMock.AssertExpectations(t)
		})
	}
}

func TestService_GetReceiver_Slack(t *testing.T) {
	var (
		timeNow    = time.Now()
		validToken = "xxxxxx"
	)

	type testCase struct {
		Description string
		Setup       func(*mocks.SecureServiceProxy, *mocks.SlackClient, testCase)
		ID          uint64
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should call repository Get method and return result in domain's type",
			Setup: func(ss *mocks.SecureServiceProxy, sc *mocks.SlackClient, tc testCase) {
				ss.EXPECT().GetReceiver(tc.ID).Return(tc.Rcv, nil)
				sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("slack.ClientCallOption"), mock.AnythingOfType("slack.ClientCallOption")).Return([]slack.Channel{
					{ID: "1", Name: "foo"},
				}, nil)
			},
			ID: 10,
			Rcv: &receiver.Receiver{
				Id:   10,
				Name: "foo",
				Type: "slack",
				Labels: map[string]string{
					"foo": "bar",
				},
				Configurations: map[string]interface{}{
					"token": validToken,
				},
				Data: map[string]interface{}{
					"channels": "[{\"id\":\"1\",\"name\":\"foo\"}]",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			Err: nil,
		},
		{
			Description: "should call secure service Get method and return error if error",
			Setup: func(ss *mocks.SecureServiceProxy, sc *mocks.SlackClient, tc testCase) {
				ss.EXPECT().GetReceiver(tc.ID).Return(nil, tc.Err)
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error if no token field in configurations",
			Setup: func(ss *mocks.SecureServiceProxy, sc *mocks.SlackClient, tc testCase) {
				ss.EXPECT().GetReceiver(tc.ID).Return(&receiver.Receiver{
					Id:   123,
					Type: "slack",
				}, nil)
			},
			Err: errors.New("no token found in configurations"),
		},
		{
			Description: "should return error if failed to get workspace channels with slack client",
			Setup: func(ss *mocks.SecureServiceProxy, sc *mocks.SlackClient, tc testCase) {
				ss.EXPECT().GetReceiver(tc.ID).Return(&receiver.Receiver{
					Id:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": validToken,
					},
				}, nil)
				sc.EXPECT().GetWorkspaceChannels(mock.AnythingOfType("slack.ClientCallOption"), mock.AnythingOfType("slack.ClientCallOption")).Return(nil, errors.New("some error"))
			},
			Err: errors.New("could not get channels: some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				secureServiceMock = new(mocks.SecureServiceProxy)
				slackClientMock   = new(mocks.SlackClient)
			)
			svc := receiver.NewService(secureServiceMock, slackClientMock)
			tc.Setup(secureServiceMock, slackClientMock, tc)

			got, err := svc.GetReceiver(tc.ID)
			if tc.Err != err {
				if tc.Err != nil && tc.Err.Error() != err.Error() {
					t.Fatalf("error %+v expected was %+v", err, tc.Err)
				}
			}
			if !cmp.Equal(got, tc.Rcv) {
				t.Fatalf("got result %+v expected was %+v", got, tc.Rcv)
			}

			secureServiceMock.AssertExpectations(t)
			slackClientMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateReceiver_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.SecureServiceProxy, testCase)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should call service Update method and return nil error",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().UpdateReceiver(tc.Rcv).Return(nil)
			},
			Rcv: &receiver.Receiver{
				Id: 123,
			},
			Err: nil,
		},
		{
			Description: "should call service update method and return error if error",
			Setup: func(ss *mocks.SecureServiceProxy, tc testCase) {
				ss.EXPECT().UpdateReceiver(tc.Rcv).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
			},
			Err: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		var (
			secureServiceMock = new(mocks.SecureServiceProxy)
			slackClientMock   = new(mocks.SlackClient)
		)
		svc := receiver.NewService(secureServiceMock, slackClientMock)

		tc.Setup(secureServiceMock, tc)

		err := svc.UpdateReceiver(tc.Rcv)

		if tc.Err != err {
			if tc.Err != nil && tc.Err.Error() != err.Error() {
				t.Fatalf("error %+v expected was %+v", err, tc.Err)
			}
		}

		secureServiceMock.AssertExpectations(t)
		slackClientMock.AssertExpectations(t)
	}
}
