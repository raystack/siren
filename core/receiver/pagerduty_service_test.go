package receiver_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
)

func TestPagerDutyService_Functions(t *testing.T) {
	t.Run("should return error not implemented if notify is called", func(t *testing.T) {
		svc := receiver.NewPagerDutyService()

		expectedErrorString := "operation not supported"
		err := svc.Notify(context.TODO(), &receiver.Receiver{}, receiver.NotificationMessage{})

		if err.Error() != expectedErrorString {
			t.Fatalf("got error %s, expected was %s", err.Error(), expectedErrorString)
		}
	})

	t.Run("should return error nil if encrypt is called", func(t *testing.T) {
		svc := receiver.NewPagerDutyService()

		err := svc.Encrypt(&receiver.Receiver{})

		if err != nil {
			t.Fatalf("got error %s, expected was nil", err.Error())
		}
	})

	t.Run("should return error nil if decrypt is called", func(t *testing.T) {
		svc := receiver.NewPagerDutyService()

		err := svc.Decrypt(&receiver.Receiver{})

		if err != nil {
			t.Fatalf("got error %s, expected was nil", err.Error())
		}
	})

	t.Run("should return as-is if populate receiver is called", func(t *testing.T) {
		svc := receiver.NewPagerDutyService()

		inputReceiver := &receiver.Receiver{
			ID:   123,
			Name: "a-receiver",
		}
		got, err := svc.PopulateReceiver(context.TODO(), inputReceiver)

		if err != nil {
			t.Fatalf("got error %s, expected was nil", err.Error())
		}

		if !cmp.Equal(got, inputReceiver) {
			t.Fatalf("got result %v, expected was %v", got, inputReceiver)
		}
	})
}

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
