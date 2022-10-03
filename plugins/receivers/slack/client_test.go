package slack_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/plugins/receivers/slack"
	"github.com/odpf/siren/plugins/receivers/slack/mocks"
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

	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "return error when goslack client creation error",
				Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
					return c.GetWorkspaceChannels(ctx)
				},
				Err: errors.New("goslack client creation failure: no client id/secret credential provided"),
			},
			{
				Description: "return error when failed to fetch joined channel list",
				Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
					gsc.EXPECT().GetConversationsForUserContext(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, "", errors.New("some error"))
					return c.GetWorkspaceChannels(ctx, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("failed to fetch joined channel list: some error"),
			},
			{
				Description: "return channels when GetWorkspaceChannels succeed",
				Call: func(c *slack.Client, gsc *mocks.GoSlackCaller) ([]slack.Channel, error) {
					gsc.EXPECT().GetConversationsForUserContext(mock.AnythingOfType("*context.emptyCtx"), &goslack.GetConversationsForUserParameters{
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
					return c.GetWorkspaceChannels(ctx, slack.CallWithGoSlackClient(gsc))
				},
				Channels: []slack.Channel{{
					ID:   "123",
					Name: "test",
				}},
				Err: nil,
			},
		}
	)
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
		Description        string
		SusbcriptionConfig slack.SubscriptionConfig
		Call               func(*slack.Client, slack.SubscriptionConfig, *mocks.GoSlackCaller) error
		Err                error
	}

	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "return error when message receiver type is wrong",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: "random",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("unknown receiver type \"random\""),
			},
			{
				Description: "(channels) return error when goslack client creation error",
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{})
				},
				Err: errors.New("goslack client creation failure: no client id/secret credential provided"),
			},
			{
				Description: "(channel) return error when failed to fetch joined channel list",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetConversationsForUserContext(mock.AnythingOfType("*context.emptyCtx"), mock.Anything).Return(nil, "", errors.New("some error"))
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("failed to fetch joined channel list: some error"),
			},
			{
				Description: "(channel) return error when app is not part of the channel",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
					ChannelName: "unknown",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetConversationsForUserContext(mock.AnythingOfType("*context.emptyCtx"), &goslack.GetConversationsForUserParameters{
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
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("app is not part of the channel \"unknown\""),
			},
			{
				Description: "(user) return error when failed to get user for an email",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetUserByEmailContext(mock.AnythingOfType("*context.emptyCtx"), "email@email.com").Return(nil, errors.New("users_not_found"))
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("failed to get id for \"email@email.com\""),
			},
			{
				Description: "(user) return error when GetUserByEmailContext return error",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetUserByEmailContext(mock.AnythingOfType("*context.emptyCtx"), "email@email.com").Return(nil, errors.New("some error"))
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "return nil error when notify is succeed through channel",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelChannel,
					ChannelName: "unknown",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetConversationsForUserContext(mock.AnythingOfType("*context.emptyCtx"), &goslack.GetConversationsForUserParameters{
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
					gsc.EXPECT().SendMessageContext(mock.AnythingOfType("*context.emptyCtx"), "123123", mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
					return c.Notify(context.TODO(), slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
			},
			{
				Description: "return nil error when notify is succeed through user",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetUserByEmailContext(mock.AnythingOfType("*context.emptyCtx"), "email@email.com").Return(&goslack.User{
						ID:   "123123",
						Name: "email@email.com",
					}, nil)
					gsc.EXPECT().SendMessageContext(mock.AnythingOfType("*context.emptyCtx"), "123123", mock.AnythingOfType("slack.MsgOption")).Return("", "", "", nil)
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
			},
			{
				Description: "return error when send message return error",
				SusbcriptionConfig: slack.SubscriptionConfig{
					ChannelType: slack.TypeChannelUser,
					ChannelName: "email@email.com",
				},
				Call: func(c *slack.Client, sc slack.SubscriptionConfig, gsc *mocks.GoSlackCaller) error {
					gsc.EXPECT().GetUserByEmailContext(mock.AnythingOfType("*context.emptyCtx"), "email@email.com").Return(&goslack.User{
						ID:   "123123",
						Name: "email@email.com",
					}, nil)
					gsc.EXPECT().SendMessageContext(mock.AnythingOfType("*context.emptyCtx"), "123123", mock.AnythingOfType("slack.MsgOption")).Return("", "", "", errors.New("some error"))
					return c.Notify(ctx, slack.NotificationConfig{SubscriptionConfig: sc}, slack.Message{}, slack.CallWithGoSlackClient(gsc))
				},
				Err: errors.New("failed to send message to \"email@email.com\""),
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			mockGoSlackClient := new(mocks.GoSlackCaller)
			c := slack.NewClient()
			err := tc.Call(c, tc.SusbcriptionConfig, mockGoSlackClient)
			if err != tc.Err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %+v, expected was %+v", err, tc.Err)
				}
			}
		})
	}
}
