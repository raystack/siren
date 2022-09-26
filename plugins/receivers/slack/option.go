package slack

import (
	"net/http"
)

type ClientOption func(*Client)

// ClientWithHTTPClient assigns custom http client when creating a slack client
func ClientWithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

type ClientCallOption func(*clientData)

// CallWithGoSlackClient uses a custom slack client when calling slack API
func CallWithGoSlackClient(gsc GoSlackCaller) ClientCallOption {
	return func(c *clientData) {
		c.goslackClient = gsc
	}
}

// CallWithClientSecret uses a new client with client ID, secret, and auth code to call slack API
func CallWithClientSecret(authCode string, clientID, clientSecret string) ClientCallOption {
	return func(c *clientData) {
		c.authCode = authCode
		c.clientID = clientID
		c.clientSecret = clientSecret
	}
}

// CallWithToken uses access token to call slack API
func CallWithToken(token string) ClientCallOption {
	return func(c *clientData) {
		c.token = token
	}
}
