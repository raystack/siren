package pagerduty_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/pkg/secret"
	"github.com/odpf/siren/plugins/receivers/pagerduty"
)

func TestPagerDutyService_Functions(t *testing.T) {
	t.Run("should return error not implemented if notify is called", func(t *testing.T) {
		svc := pagerduty.NewReceiverService()

		expectedErrorString := "operation not supported"
		err := svc.Notify(context.TODO(), map[string]interface{}{}, map[string]interface{}{})

		if err.Error() != expectedErrorString {
			t.Fatalf("got error %s, expected was %s", err.Error(), expectedErrorString)
		}
	})

	t.Run("should return empty if get populated data is called", func(t *testing.T) {
		svc := pagerduty.NewReceiverService()

		got, err := svc.BuildData(context.TODO(), map[string]interface{}{})
		if err != nil {
			t.Fatalf("got error %s, expected was nil", err.Error())
		}

		if !cmp.Equal(got, map[string]interface{}{}) {
			t.Fatalf("got result %v, expected was %v", got, map[string]interface{}{})
		}
	})
}

func TestPagerDutyService_BuildNotificationConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]interface{}
		ReceiverConfigs     map[string]interface{}
		ExpectedConfigMap   map[string]interface{}
		wantErr             bool
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if receiver 'service_key' exist but it is not string",
				ReceiverConfigs: map[string]interface{}{
					"service_key": 123,
				},
				wantErr: true,
			},
			{
				Description: "should return configs with service_key if receiver 'service_key' exist in string",
				ReceiverConfigs: map[string]interface{}{
					"service_key": secret.MaskableString("service_key"),
				},
				ExpectedConfigMap: map[string]interface{}{
					"service_key": secret.MaskableString("service_key"),
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := pagerduty.NewReceiverService()

			got, err := svc.BuildNotificationConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
			if (err != nil) != tc.wantErr {
				t.Errorf("got error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil {
				if !cmp.Equal(got, tc.ExpectedConfigMap) {
					t.Errorf("got result %+v, expected was %+v", got, tc.ExpectedConfigMap)
				}
			}
		})
	}
}
