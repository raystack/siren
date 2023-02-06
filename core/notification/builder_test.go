package notification_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/template"
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

func TestBuildFromAlerts(t *testing.T) {
	tests := []struct {
		name      string
		alerts    []alert.Alert
		firingLen int
		want      notification.Notification
	}{

		{
			name: "should return empty notification if alerts slice is empty",
			want: notification.Notification{},
		},
		{
			name: `should properly return notification
				- same annotations are joined by newline
				- labels are merged
			`,
			alerts: []alert.Alert{
				{
					ID:           14,
					ProviderID:   1,
					NamespaceID:  1,
					ResourceName: "test-alert-host-1",
					MetricName:   "test-alert",
					MetricValue:  "15",
					Severity:     "WARNING",
					Rule:         "test-alert-template",
					Labels:       map[string]string{"lk1": "lv1"},
					Annotations:  map[string]string{"ak1": "akv1"},
					Status:       "FIRING",
				},
				{
					ID:           15,
					ProviderID:   1,
					NamespaceID:  1,
					ResourceName: "test-alert-host-2",
					MetricName:   "test-alert",
					MetricValue:  "16",
					Severity:     "WARNING",
					Rule:         "test-alert-template",
					Labels:       map[string]string{"lk1": "lv11", "lk2": "lv2"},
					Annotations:  map[string]string{"ak1": "akv11", "ak2": "akv2"},
					Status:       "FIRING",
				},
			},
			firingLen: 2,
			want: notification.Notification{
				NamespaceID: 1,
				Type:        notification.TypeSubscriber,

				Data: map[string]interface{}{
					"generator_url":     "",
					"num_alerts_firing": 2,
					"status":            "FIRING",
					"ak1":               "akv1\nakv11",
					"ak2":               "akv2",
				},
				Labels: map[string]string{
					"lk1": "lv11",
					"lk2": "lv2",
				},
				Template: template.ReservedName_SystemDefault,
				AlertIDs: []int64{14, 15},
			},
		},
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := notification.BuildFromAlerts(tt.alerts, tt.firingLen, time.Time{})

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("BuildFromAlerts() got diff = %v", diff)
			}
		})
	}
}
