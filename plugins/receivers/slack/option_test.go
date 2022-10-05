package slack_test

import (
	"testing"

	"github.com/odpf/siren/pkg/httpclient"
	"github.com/odpf/siren/plugins/receivers/slack"
)

func TestClientOption(t *testing.T) {
	t.Run("should be fine when using http client from external", func(t *testing.T) {
		c := slack.NewClient(slack.AppConfig{}, slack.ClientWithHTTPClient(httpclient.New(httpclient.Config{})))
		if c == nil {
			t.Fatal("client should not be nil")
		}
	})
}

// func TestClientCallOption(t *testing.T) {
// 	t.Run("should use token when token is passed", func(t *testing.T) {
// 		c := slack.NewClient()
// 		if c == nil {
// 			t.Fatal("client should not be nil")
// 		}

// 		_, err := c.GetWorkspaceChannels(
// 			context.TODO(),
// 			slack.CallWithToken("1234"),
// 		)

// 		expectedErrorString := "failed to fetch joined channel list: invalid_auth"
// 		if err.Error() != expectedErrorString {
// 			t.Fatalf("got error %s, expected was %s", err, expectedErrorString)
// 		}
// 	})

// 	t.Run("should use client id and secret when passed", func(t *testing.T) {
// 		c := slack.NewClient()
// 		if c == nil {
// 			t.Fatal("client should not be nil")
// 		}

// 		_, err := c.GetWorkspaceChannels(
// 			context.TODO(),
// 			slack.CallWithClientSecret("1234", "1234", "1234"),
// 		)

// 		expectedErrorString := "goslack client creation failure: slack oauth call failed"
// 		if err.Error() != expectedErrorString {
// 			t.Fatalf("got error %s, expected was %s", err, expectedErrorString)
// 		}
// 	})

// 	t.Run("should return error if use client id and secret but missing required params", func(t *testing.T) {
// 		c := slack.NewClient()
// 		if c == nil {
// 			t.Fatal("client should not be nil")
// 		}

// 		_, err := c.GetWorkspaceChannels(
// 			context.TODO(),
// 			slack.CallWithClientSecret("1234", "1234", ""),
// 		)

// 		expectedErrorString := "goslack client creation failure: no client id/secret credential provided"
// 		if err.Error() != expectedErrorString {
// 			t.Fatalf("got error %s, expected was %s", err, expectedErrorString)
// 		}
// 	})
// }
