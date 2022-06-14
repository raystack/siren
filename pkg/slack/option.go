package slack

import (
	"context"
	"net/http"
)

type ClientOption func(*Client)

func ClientWithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

type ClientCallOption func(*clientData)

func CallWithGoSlackClient(gsc GoSlackCaller) ClientCallOption {
	return func(c *clientData) {
		c.goslackClient = gsc
	}
}

func CallWithClientSecret(authCode string, clientID, clientSecret string) ClientCallOption {
	return func(c *clientData) {
		c.authCode = authCode
		c.clientID = clientID
		c.clientSecret = clientSecret
	}
}

func CallWithToken(token string) ClientCallOption {
	return func(c *clientData) {
		c.token = token
	}
}

func CallWithContext(ctx context.Context) ClientCallOption {
	return func(c *clientData) {
		c.ctx = ctx
	}
}
