package receiver_test

import (
	"testing"

	"github.com/odpf/siren/core/receiver"
)

func TestPagerDutyService_ValidateConfiguration(t *testing.T) {
	type testCase struct {
		Description  string
		InputConfigs receiver.Configurations
		ErrString    string
	}

	var (
		testCases = []testCase{
			{
				Description: "should return error if 'service_key' is empty",
				ErrString:   "no value supplied for required configurations map key \"service_key\"",
			},
			{
				Description: "should return nil error if all configurations are valid",
				InputConfigs: receiver.Configurations{
					"service_key": "service_key",
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			svc := receiver.NewPagerDutyService()

			err := svc.ValidateConfiguration(tc.InputConfigs)
			if err != nil {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}