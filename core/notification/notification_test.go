package notification_test

import (
	"testing"

	"github.com/odpf/siren/core/notification"
)

func TestNotification_Validate(t *testing.T) {
	testCases := []struct {
		name                string
		n                   notification.Notification
		receiverType        string
		notificationConfigs map[string]interface{}
		wantErr             bool
	}{
		{
			name: "should return error if type is unknown",
			n: notification.Notification{
				Type: "random",
			},
			wantErr: true,
		},
		{
			name: "should return error if type receiver but has no 'receiver_id' key label",
			n: notification.Notification{
				Type: notification.TypeReceiver,
				Labels: map[string]string{
					"labelkey1": "value1",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if type receiver but has empty 'receiver_id' label value",
			n: notification.Notification{
				Type: notification.TypeReceiver,
				Labels: map[string]string{
					"receiver_id": "",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return error if type receiver but has 'receiver_id' value is non parsable to integer",
			n: notification.Notification{
				Type: notification.TypeReceiver,
				Labels: map[string]string{
					"receiver_id": "xxx",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return nil error if type receiver and 'receiver_id' is valid",
			n: notification.Notification{
				Type: notification.TypeReceiver,
				Labels: map[string]string{
					"receiver_id": "2",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
		},
		{
			name: "should return error if type subscriber but has no kv labels",
			n: notification.Notification{
				Type: notification.TypeSubscriber,
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
			wantErr: true,
		},
		{
			name: "should return nil error if type subscriber and has kv labels",
			n: notification.Notification{
				Type: notification.TypeSubscriber,
				Labels: map[string]string{
					"receiver_id": "xxx",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.n.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Notification.ToMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
		})
	}
}
