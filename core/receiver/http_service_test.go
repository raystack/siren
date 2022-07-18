package receiver_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
)

func TestHTTPService_ValidateConfiguration(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		ErrString   string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if 'url' is empty",
				Rcv:         &receiver.Receiver{},
				ErrString:   "no value supplied for required configurations map key \"url\"",
			},
			{
				Description: "should return nil error if all configurations are valid",
				Rcv: &receiver.Receiver{
					Configurations: receiver.Configurations{
						"url": "url",
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
			svc := receiver.NewHTTPService()

			err := svc.ValidateConfiguration(tc.Rcv)
			if err != nil {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func TestHTTPService_GetSubscriptionConfig(t *testing.T) {
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
				Description: "should return error if receiver 'url' exist but it is not string",
				ReceiverConfigs: receiver.Configurations{
					"url": 123,
				},
				ErrString: "url config from receiver should be in string",
			},
			{
				Description:       "should return configs without token if receiver 'url' does not exist", //TODO might need to check this behaviour, should be returning error
				ExpectedConfigMap: map[string]string{},
			},
			{
				Description: "should return configs with token if receiver 'url' exist in string",
				ReceiverConfigs: receiver.Configurations{
					"url": "url",
				},
				ExpectedConfigMap: map[string]string{
					"url": "url",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := receiver.NewHTTPService()

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
