package slackchannel_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/siren/pkg/secret"
	"github.com/goto/siren/plugins/receivers/slack"
	"github.com/goto/siren/plugins/receivers/slackchannel"
)

func TestReceiverConfig(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		testCases := []struct {
			name    string
			c       slackchannel.ReceiverConfig
			wantErr bool
		}{
			{
				name:    "return error if one of required field is missing",
				wantErr: true,
			},
			{
				name: "return nil if all required fields are present",
				c: slackchannel.ReceiverConfig{
					ChannelName: "a-channel",
				},
				wantErr: false,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.c.Validate(); (err != nil) != tc.wantErr {
					t.Errorf("ReceiverConfig.Validate() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})
}

func TestNotificationConfig(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		testCases := []struct {
			name    string
			c       slackchannel.NotificationConfig
			wantErr bool
		}{
			{
				name:    "return error if one of required field is missing",
				wantErr: true,
			},
			{
				name: "return nil if all required fields are present",
				c: slackchannel.NotificationConfig{
					ReceiverConfig: slackchannel.ReceiverConfig{
						SlackReceiverConfig: slack.ReceiverConfig{
							Token:     "token",
							Workspace: "workspace",
						},
						ChannelName: "a-channel",
					},
				},
				wantErr: false,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.c.Validate(); (err != nil) != tc.wantErr {
					t.Errorf("NotificationConfig.Validate() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})

	t.Run("AsMap", func(t *testing.T) {
		nc := slackchannel.NotificationConfig{
			ReceiverConfig: slackchannel.ReceiverConfig{
				SlackReceiverConfig: slack.ReceiverConfig{
					Token:     secret.MaskableString("token"),
					Workspace: "workspace",
				},
				ChannelName: "channel",
			},
		}

		if diff := cmp.Diff(map[string]any{
			"channel_name": "channel",
			"channel_type": "",
			"token":        secret.MaskableString("token"),
			"workspace":    "workspace",
		}, nc.AsMap()); diff != "" {
			t.Errorf("result not match\n%v", diff)
		}
	})
}
