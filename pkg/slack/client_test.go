package slack_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/pkg/slack"
	"github.com/odpf/siren/pkg/slack/mocks"
	goslack "github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetWorkspaceChannels(t *testing.T) {
	type testCase struct {
		Description string
		Call        func(*slack.Client, *mocks.GoSlackCaller) ([]slack.Channel, error)
		Channels    []slack.Channel
		Err         error
	}

	var testCases = []testCase{
		{
			Description: "return error when goslack client creation error",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
				return c.GetWorkspaceChannels()
			},
			Err: errors.New("goslack client creation failure: no client id/secret credential provided"),
		},
		{
			Description: "return error when failed to fetch joined channel list",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
				gsc.EXPECT().GetConversationsForUser(mock.Anything).Return(nil, "", errors.New("some error"))
				return c.GetWorkspaceChannels(slack.CallWithGoSlackClient(gsc))
			},
			Err: errors.New("failed to fetch joined channel list: some error"),
		},
		{
			Description: "return channels when GetWorkspaceChannels succeed",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
				gsc.EXPECT().GetConversationsForUser(&goslack.GetConversationsForUserParameters{
					Types:  []string{"public_channel", "private_channel"},
					Cursor: "",
					Limit:  1000}).Return([]goslack.Channel{
					{
						GroupConversation: goslack.GroupConversation{
							Conversation: goslack.Conversation{
								ID: "123",
							},
							Name: "test",
						},
						IsChannel: true,
					},
				}, "", nil)
				return c.GetWorkspaceChannels(slack.CallWithGoSlackClient(gsc))
			},
			Channels: []slack.Channel{{
				ID:   "123",
				Name: "test",
			}},
			Err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			mockGoSlackClient := new(mocks.GoSlackCaller)
			c := slack.NewClient()

			got, err := tc.Call(c, mockGoSlackClient)
			if err != tc.Err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %+v, expected was %+v", err, tc.Err)
				}
			}
			if !cmp.Equal(got, tc.Channels) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.Channels)
			}
		})
	}
}

func TestClient_Notify(t *testing.T) {
	type testCase struct {
		Description string
		Call        func(*slack.Client, *mocks.GoSlackCaller) error
		Err         error
	}

	var testCases = []testCase{
		{
			Description: "return error when message receiver type is wrong",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				return c.Notify(&slack.Message{
					ReceiverType: "random",
				}, slack.CallWithGoSlackClient(gsc))
			},
			Err: errors.New("unknown receiver type 'random'"),
		},
		{
			Description: "(channels) return error when goslack client creation error",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				return c.Notify(nil)
			},
			Err: errors.New("goslack client creation failure: no client id/secret credential provided"),
		},
		{
			Description: "(channel) return error when failed to fetch joined channel list",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				gsc.EXPECT().GetConversationsForUser(mock.Anything).Return(nil, "", errors.New("some error"))
				return c.Notify(&slack.Message{
					ReceiverType: "channel", //TODO make it const
				}, slack.CallWithGoSlackClient(gsc))
			},
			Err: errors.New("failed to fetch joined channel list: some error"),
		},
		{
			Description: "(channel) return error when app is not part of the channel",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				gsc.EXPECT().GetConversationsForUser(&goslack.GetConversationsForUserParameters{
					Types:  []string{"public_channel", "private_channel"},
					Cursor: "",
					Limit:  1000}).Return([]goslack.Channel{
					{
						GroupConversation: goslack.GroupConversation{
							Conversation: goslack.Conversation{
								ID: "123",
							},
							Name: "test",
						},
						IsChannel: true,
					},
				}, "", nil)
				return c.Notify(&slack.Message{
					ReceiverName: "unknwon",
					ReceiverType: "channel", //TODO make it const
				}, slack.CallWithGoSlackClient(gsc))
			},
			Err: errors.New("app is not part of the channel unknwon"),
		},
		{
			Description: "(user) return error when failed to get user for an email",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				gsc.EXPECT().GetUserByEmail("email@email.com").Return(nil, errors.New("users_not_found"))
				return c.Notify(&slack.Message{
					ReceiverName: "email@email.com",
					ReceiverType: "user", //TODO make it const
				}, slack.CallWithGoSlackClient(gsc))
			},
			Err: errors.New("failed to get id for email@email.com"),
		},
		{
			Description: "return nil error when notify is succeed",
			Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) error {
				gsc.EXPECT().GetConversationsForUser(&goslack.GetConversationsForUserParameters{
					Types:  []string{"public_channel", "private_channel"},
					Cursor: "",
					Limit:  1000}).Return([]goslack.Channel{
					{
						GroupConversation: goslack.GroupConversation{
							Conversation: goslack.Conversation{
								ID: "123123",
							},
							Name: "unknown",
						},
						IsChannel: true,
					},
				}, "", nil)
				gsc.EXPECT().SendMessage("123123", mock.AnythingOfType("slack.MsgOption"), mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
				return c.Notify(&slack.Message{
					ReceiverName: "unknown",
					ReceiverType: "channel", //TODO make it const
				}, slack.CallWithGoSlackClient(gsc))
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			mockGoSlackClient := new(mocks.GoSlackCaller)
			c := slack.NewClient()
			err := tc.Call(c, mockGoSlackClient)
			if err != tc.Err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %+v, expected was %+v", err, tc.Err)
				}
			}
		})
	}
}
