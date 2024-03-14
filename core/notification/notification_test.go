package notification_test

import (
	"testing"

	"github.com/goto/siren/core/notification"
)

func TestNotification_Validate(t *testing.T) {
	testCases := []struct {
		name                string
		n                   notification.Notification
		Flow                string
		receiverType        string
		notificationConfigs map[string]any
		wantErr             bool
	}{
		{
			name:    "should return error if flow is unknown",
			Flow:    "random",
			n:       notification.Notification{},
			wantErr: true,
		},
		{
			name: "should return error if flow receiver but has no receiver_selectors",
			Flow: notification.FlowReceiver,
			n: notification.Notification{
				Labels: map[string]string{
					"labelkey1": "value1",
				},
				Data: map[string]any{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return nil error if flow receiver and receiver_selectors exist",
			Flow: notification.FlowReceiver,
			n: notification.Notification{
				Labels: map[string]string{
					"receiver_id": "2",
				},
				ReceiverSelectors: []map[string]string{
					{
						"varkey1": "value1",
					},
				},
			},
		},
		{
			name: "should return error if flow subscriber but has no kv labels",
			Flow: notification.FlowSubscriber,
			n: notification.Notification{
				Data: map[string]any{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return nil error if flow subscriber and has kv labels",
			Flow: notification.FlowSubscriber,
			n: notification.Notification{
				Labels: map[string]string{
					"receiver_id": "xxx",
				},
				Data: map[string]any{
					"varkey1": "value1",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.n.Validate(tc.Flow)
			if (err != nil) != tc.wantErr {
				t.Errorf("Notification.ToMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
