package slack_test

import (
	"testing"

	"github.com/odpf/siren/pkg/slack"
)

func TestMessage_Validate(t *testing.T) {
	type testCase struct {
		Description       string
		Message           slack.Message
		ExpectedErrString string
	}

	var testCases = []testCase{
		{
			Description:       "should return error message or blocks cannot be empty",
			Message:           slack.Message{},
			ExpectedErrString: "non empty message or non zero length block is required",
		},
		{
			Description: "should return error if required fields are not populated",
			Message: slack.Message{
				Message: "a message",
			},
			ExpectedErrString: "field \"receiver_name\" is required and field \"receiver_type\" is required",
		},
		{
			Description: "should return error type not supported if slack receiver type not match",
			Message: slack.Message{
				Message:      "a message",
				ReceiverName: "receiver name",
				ReceiverType: "random",
			},
			ExpectedErrString: "error value \"random\" for key \"receiver_type\" not recognized, only support \"user channel\"",
		},
		{
			Description: "should return multiple validation errors",
			Message: slack.Message{
				Message:      "a message",
				ReceiverType: "random",
			},
			ExpectedErrString: "field \"receiver_name\" is required and error value \"random\" for key \"receiver_type\" not recognized, only support \"user channel\"",
		},
		{
			Description: "should return nil error if slack message is valid",
			Message: slack.Message{
				Message:      "a message",
				ReceiverName: "receiver name",
				ReceiverType: "channel",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			err := tc.Message.Validate()
			if err != nil {
				if err.Error() != tc.ExpectedErrString {
					t.Fatalf("got error %q, expected was %q", err.Error(), tc.ExpectedErrString)
				}
			}
		})
	}
}
