package cortex_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/plugins/providers/cortex"
)

func TestGetAlertManagerReceiverConfig(t *testing.T) {
	type testCase struct {
		Description              string
		Subscription             *subscription.Subscription
		ExpectedAMReceiverConfig []cortex.ReceiverConfig
	}

	var testCases = []testCase{
		{
			Description: "should return nil if subscription is nil",
		},
		{
			Description: "should build am receiver configs properly",
			Subscription: &subscription.Subscription{
				Receivers: []subscription.Receiver{
					{
						ID:   5,
						Type: receiver.TypeHTTP,
						Configuration: map[string]interface{}{
							"url": "http://webhook",
						},
					},
					{
						ID:   7,
						Type: receiver.TypeSlack,
						Configuration: map[string]interface{}{
							"channel_name": "odpf-channel",
							"token":        "123123123",
						},
					},
					{
						ID:   9,
						Type: receiver.TypeSlack,
						Configuration: map[string]interface{}{
							"channel_name": "odpf-channel",
							"token":        "123123123",
						},
					},
					{
						ID:   10,
						Type: receiver.TypePagerDuty,
						Configuration: map[string]interface{}{
							"service_key": "a-service-key",
						},
					},
				},
				Match: map[string]string{
					"label1": "value1",
					"label2": "value2",
					"label3": "value3",
				},
			},
			ExpectedAMReceiverConfig: []cortex.ReceiverConfig{
				{
					Name: "_receiverId_5_idx_0",
					Type: receiver.TypeHTTP,
					Match: map[string]string{
						"label1": "value1",
						"label2": "value2",
						"label3": "value3",
					},
					Configurations: map[string]string{
						"url": "http://webhook",
					},
				},
				{
					Name: "_receiverId_7_idx_1",
					Type: receiver.TypeSlack,

					Match: map[string]string{
						"label1": "value1",
						"label2": "value2",
						"label3": "value3",
					},
					Configurations: map[string]string{
						"channel_name": "odpf-channel",
						"token":        "123123123",
					},
				},
				{
					Name: "_receiverId_9_idx_2",
					Type: receiver.TypeSlack,
					Match: map[string]string{
						"label1": "value1",
						"label2": "value2",
						"label3": "value3",
					},
					Configurations: map[string]string{
						"channel_name": "odpf-channel",
						"token":        "123123123",
					},
				},
				{
					Name: "_receiverId_10_idx_3",
					Type: receiver.TypePagerDuty,
					Match: map[string]string{
						"label1": "value1",
						"label2": "value2",
						"label3": "value3",
					},
					Configurations: map[string]string{
						"service_key": "a-service-key",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			got := cortex.GetAlertManagerReceiverConfig(tc.Subscription)
			if !cmp.Equal(got, tc.ExpectedAMReceiverConfig) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedAMReceiverConfig)
			}
		})
	}

}
