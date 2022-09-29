package pagerduty_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
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

// func TestPagerDutyService_ValidateConfigurations(t *testing.T) {
// 	type testCase struct {
// 		Description string
// 		Confs       map[string]interface{}
// 		ErrString   string
// 	}

// 	var (
// 		testCases = []testCase{
// 			{
// 				Description: "should return error if 'service_key' is empty",
// 				Confs:       make(map[string]interface{}),
// 				ErrString:   "no value supplied for required configurations map key \"service_key\"",
// 			},
// 			{
// 				Description: "should return nil error if all configurations are valid",
// 				Confs: map[string]interface{}{
// 					"service_key": "service_key",
// 				},
// 			},
// 		}
// 	)

// 	for _, tc := range testCases {
// 		t.Run(tc.Description, func(t *testing.T) {
// 			svc := pagerduty.NewReceiverService()

// 			err := svc.ValidateConfigurations(tc.Confs)
// 			if err != nil {
// 				if tc.ErrString != err.Error() {
// 					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
// 				}
// 			}
// 		})
// 	}
// }

func TestPagerDutyService_GetSubscriptionConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]interface{}
		ReceiverConfigs     map[string]interface{}
		ExpectedConfigMap   map[string]string
		ErrString           string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if receiver 'service_key' exist but it is not string",
				ReceiverConfigs: map[string]interface{}{
					"service_key": 123,
				},
				ErrString: "service_key config from receiver should be in string",
			},
			{
				Description:       "should return configs without token if receiver 'service_key' does not exist", //TODO might need to check this behaviour, should be returning error
				ExpectedConfigMap: map[string]string{},
			},
			{
				Description: "should return configs with token if receiver 'service_key' exist in string",
				ReceiverConfigs: map[string]interface{}{
					"service_key": "service_key",
				},
				ExpectedConfigMap: map[string]string{
					"service_key": "service_key",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := pagerduty.NewReceiverService()

			got, err := svc.BuildNotificationConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedConfigMap) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedConfigMap)
			}
		})
	}
}
