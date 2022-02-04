package codeexchange

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const OAuthServerEndpoint = "https://slack.com/api/oauth.v2.access"

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Exchanger interface {
	Exchange(string, string, string) (CodeExchangeHTTPResponse, error)
}

type SlackClient struct {
	httpClient Doer
}

func NewSlackClient(doer Doer) *SlackClient {
	return &SlackClient{httpClient: doer}
}

func (c *SlackClient) Exchange(code, clientID, clientSecret string) (CodeExchangeHTTPResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	response := CodeExchangeHTTPResponse{}
	req, err := http.NewRequest(http.MethodPost, OAuthServerEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return response, errors.Wrap(err, "failed to create http request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return response, errors.Wrap(err, "failure in http call")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, errors.Wrap(err, "failed to read response body")
	}
	err = json.Unmarshal(bodyBytes, &response)

	if err != nil {
		return response, errors.Wrap(err, "failed to unmarshal response body")
	}
	if !response.Ok {
		return response, errors.New("slack oauth call failed")
	}
	fmt.Println(response)
	return response, nil
}
