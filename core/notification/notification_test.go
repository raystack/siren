package notification_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/siren/core/notification"
)

func TestNotification_ToMessage(t *testing.T) {
	testCases := []struct {
		name                string
		n                   notification.Notification
		receiverType        string
		notificationConfigs map[string]interface{}
		want                *notification.Message
		wantErr             bool
	}{
		{
			name: "should return error if expiry duration is not parsable",
			n: notification.Notification{
				ValidDurationString: "xxx",
			},
			wantErr: true,
		},
		{
			name: "should return message if expiry duration is empty",
			n: notification.Notification{
				Labels: map[string]string{
					"labelkey1": "value1",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
			},
			want: &notification.Message{
				Status: notification.MessageStatusEnqueued,
				Detail: map[string]interface{}{
					"labelkey1": "value1",
					"varkey1":   "value1",
				},
			},
		},
		{
			name: "should return message if expiry duration is parsable",
			n: notification.Notification{
				Labels: map[string]string{
					"labelkey1": "value1",
				},
				Data: map[string]interface{}{
					"varkey1": "value1",
				},
				ValidDurationString: "10m",
			},
			want: &notification.Message{
				Status: notification.MessageStatusEnqueued,
				Detail: map[string]interface{}{
					"labelkey1": "value1",
					"varkey1":   "value1",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.n.ToMessage(tc.receiverType, tc.notificationConfigs)
			if (err != nil) != tc.wantErr {
				t.Errorf("Notification.ToMessage() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if diff := cmp.Diff(got, tc.want,
				cmpopts.IgnoreUnexported(notification.Message{}),
				cmpopts.IgnoreFields(
					notification.Message{},
					"ID",
					"MaxTries",
					"ExpiredAt",
					"CreatedAt",
					"UpdatedAt",
				),
			); diff != "" {
				t.Errorf("Notification.ToMessage() diff = %v", diff)
			}
		})
	}
}
