package notification_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/siren/core/notification"
)

func TestBuildTypeReceiver(t *testing.T) {
	sampleReceiverID := uint64(11)

	tests := []struct {
		name       string
		receiverID uint64
		payloadMap map[string]interface{}
		want       notification.Notification
		wantErr    bool
	}{
		{
			name:       "should build a correct notification",
			receiverID: sampleReceiverID,
			payloadMap: map[string]interface{}{
				"data": map[string]interface{}{
					"key1": "key2",
				},
				"valid_duration": "10m",
				"template":       "some-template",
			},
			want: notification.Notification{
				Type: notification.TypeReceiver,
				Data: map[string]interface{}{
					"key1": "key2",
				},
				Labels: map[string]string{
					"receiver_id": "11",
				},
				ValidDuration: time.Duration(10 * time.Minute),
				Template:      "some-template",
			},
		},
		{
			name:       "should return error if payload is not decodable",
			receiverID: sampleReceiverID,
			payloadMap: map[string]interface{}{
				"template": 1,
			},
			wantErr: true,
		},
		{
			name:       "should return error if 'valid_duration' is not string",
			receiverID: sampleReceiverID,
			payloadMap: map[string]interface{}{
				"valid_duration": 1,
			},
			wantErr: true,
		},
		{
			name:       "should return error if 'valid_duration' is not parsable",
			receiverID: sampleReceiverID,
			payloadMap: map[string]interface{}{
				"valid_duration": "xzx",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := notification.BuildTypeReceiver(tt.receiverID, tt.payloadMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildTypeReceiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("BuildTypeReceiver() diff = %v", diff)
			}
		})
	}
}
