package receiver_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
)

func TestPagerDutyService_ValidateConfiguration(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		ErrString   string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if 'service_key' is empty",
				Rcv:         &receiver.Receiver{},
				ErrString:   "no value supplied for required configurations map key \"service_key\"",
			},
			{
				Description: "should return nil error if all configurations are valid",
				Rcv: &receiver.Receiver{
					Configurations: receiver.Configurations{
						"service_key": "service_key",
					},
				},
			},
			{
				Description: "should return error if receiver is nil",
				ErrString:   "receiver to validate is nil",
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := receiver.NewPagerDutyService()

			err := svc.ValidateConfiguration(tc.Rcv)
			if err != nil {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestPagerDutyService_GetSubscriptionConfig(t *testing.T) {
	type testCase struct {
		Description         string
		SubscriptionConfigs map[string]string
		ReceiverConfigs     receiver.Configurations
		ExpectedConfigMap   map[string]string
		ErrString           string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if receiver 'service_key' exist but it is not string",
				ReceiverConfigs: receiver.Configurations{
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
				ReceiverConfigs: receiver.Configurations{
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
			svc := receiver.NewPagerDutyService()

			got, err := svc.GetSubscriptionConfig(tc.SubscriptionConfigs, tc.ReceiverConfigs)
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
