package slack

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSlackCredentialConfig(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		testCases := []struct {
			name    string
			c       SlackCredentialConfig
			wantErr bool
		}{
			{
				name:    "return error if one of required field is missing",
				wantErr: true,
			},
			{
				name: "return nil if all required fields are present",
				c: SlackCredentialConfig{
					ClientID:     "clientid",
					ClientSecret: "clientsecret",
					AuthCode:     "authcode",
				},
				wantErr: false,
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if err := tc.c.Validate(); (err != nil) != tc.wantErr {
					t.Errorf("RegisterReceiverConfig.Validate() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})
}

func TestReceiverConfig(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		testCases := []struct {
			name    string
			c       ReceiverConfig
			wantErr bool
		}{
			{
				name:    "return error if one of required field is missing",
				wantErr: true,
			},
			{
				name: "return nil if all required fields are present",
				c: ReceiverConfig{
					Token:     "token",
					Workspace: "workspace",
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
			c       NotificationConfig
			wantErr bool
		}{
			{
				name:    "return error if one of required field is missing",
				wantErr: true,
			},
			{
				name: "return nil if all required fields are present",
				c: NotificationConfig{
					SubscriptionConfig: SubscriptionConfig{
						ChannelName: "channel",
					},
					ReceiverConfig: ReceiverConfig{
						Token:     "token",
						Workspace: "workspace"},
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
		nc := NotificationConfig{
			SubscriptionConfig: SubscriptionConfig{
				ChannelName: "channel",
			},
			ReceiverConfig: ReceiverConfig{
				Token:     "token",
				Workspace: "workspace"},
		}

		if diff := cmp.Diff(map[string]interface{}{
			"channel_name": "channel",
			"token":        "token",
			"workspace":    "workspace",
		}, nc.AsMap()); diff != "" {
			t.Errorf("result not match\n%v", diff)
		}
	})
}
