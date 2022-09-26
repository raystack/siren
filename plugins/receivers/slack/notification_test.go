package slack_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/plugins/receivers/slack"
	goslack "github.com/slack-go/slack"
)

func TestNotification_ToSlackMessage(t *testing.T) {
	type testCase struct {
		Description          string
		Message              map[string]interface{}
		ExpectedSlackMessage *slack.Message
		ExpectedErrString    string
	}

	var testCases = []testCase{
		{
			Description: "should return error if struct can't be marshalled",
			Message: map[string]interface{}{
				"test": make(chan int),
			},
			ExpectedErrString: "unable to marshal notification message: json: unsupported type: chan int",
		},
		{
			Description: "should return error if json byte can't be unmarshalled",
			Message: map[string]interface{}{
				"blocks": "abc",
			},
			ExpectedErrString: "unable to unmarshal notification message byte to slack message: json: cannot unmarshal string into Go struct field Message.blocks of type []json.RawMessage",
		},
		{
			Description:       "should return error if 'message' are empty and blocks are empty",
			Message:           map[string]interface{}{},
			ExpectedErrString: "non empty message or non zero length block is required",
		},
		{
			Description: "should return error if required fields are empty",
			Message: map[string]interface{}{
				"blocks": []map[string]interface{}{
					{
						"key": "value",
					},
				},
			},
			ExpectedErrString: "field \"receiver_name\" is required and field \"receiver_type\" is required",
		},
		{
			Description: "should return slack message if notification message is valid",
			Message: map[string]interface{}{
				"receiver_name": "receiver_name",
				"receiver_type": "channel",
				"message":       "message",
				"blocks": []map[string]interface{}{
					{
						"type": "section",
					},
				},
			},
			ExpectedSlackMessage: &slack.Message{
				ReceiverName: "receiver_name",
				ReceiverType: "channel",
				Message:      "message",
				Blocks: goslack.Blocks{
					BlockSet: []goslack.Block{
						&goslack.SectionBlock{
							Type: goslack.MBTSection,
						},
					},
				},
			},
			ExpectedErrString: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			got, err := slack.GetSlackMessage(tc.Message)
			if err != nil {
				if err.Error() != tc.ExpectedErrString {
					t.Fatalf("got error '%s', expected was '%s'", err.Error(), tc.ExpectedErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedSlackMessage) {
				t.Fatalf("got result '%+v', expected was '%+v'", got, tc.ExpectedSlackMessage)
			}

		})
	}
}
