package codeexchange

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type HTTPTestSuite struct {
	suite.Suite
	doerMock DoerMock
}

func TestHTTP(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

func (s *HTTPTestSuite) SetupTest() {
	s.doerMock = DoerMock{}
}

func (s *HTTPTestSuite) TestHTTP_Exchange() {
	slackUrl := "https://slack.com/api/oauth.v2.access"
	clientID, clientSecret, code := "test-client-id", "test-client-secret", "test-code"

	s.Run("should call slack oauth endpoint for code exchange", func() {
		doerMock := &DoerMock{}
		dummySlackClient := SlackClient{
			httpClient: doerMock,
		}
		data := url.Values{}
		data.Set("ok", "true")
		doerMock.On("Do", mock.AnythingOfType("*http.Request")).Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType((*http.Request)(nil), rarg)
			r := rarg.(*http.Request)
			s.Equal(http.MethodPost, r.Method)
			s.Equal(slackUrl, r.URL.String())
			s.Equal("application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
			body, _ := ioutil.ReadAll(r.Body)
			s.Equal("client_id=test-client-id&client_secret=test-client-secret&code=test-code", string(body))
		}).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok":true,"access_token":"foo","team":{"name":"bar"}}`)),
		}, nil)

		resp, err := dummySlackClient.Exchange(code, clientID, clientSecret)

		s.Nil(err)
		s.Equal("foo", resp.AccessToken)
		s.Equal("bar", resp.Team.Name)
		s.Equal(true, resp.Ok)
	})

	s.Run("should handle http errors in code exchange", func() {
		doerMock := &DoerMock{}
		dummySlackClient := SlackClient{
			httpClient: doerMock,
		}
		doerMock.On("Do", mock.AnythingOfType("*http.Request")).Return(nil,
			errors.New("random error"))

		resp, err := dummySlackClient.Exchange(code, clientID, clientSecret)

		s.EqualError(err, "failure in http call: random error")
		s.NotNil(resp)
	})

	s.Run("should handle slack errors in code exchange", func() {
		doerMock := &DoerMock{}
		dummySlackClient := SlackClient{
			httpClient: doerMock,
		}
		doerMock.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"ok":false}`)),
		}, nil)

		resp, err := dummySlackClient.Exchange(code, clientID, clientSecret)

		s.EqualError(err, `slack oauth call failed`)
		s.Equal(false, resp.Ok)
	})

	s.Run("should handle http response parse errors in code exchange", func() {
		doerMock := &DoerMock{}
		dummySlackClient := SlackClient{
			httpClient: doerMock,
		}
		data := url.Values{}
		data.Set("ok", "true")
		doerMock.On("Do", mock.AnythingOfType("*http.Request")).Return(&http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(`Hello World`)),
		}, nil)

		resp, err := dummySlackClient.Exchange(code, clientID, clientSecret)

		s.EqualError(err, `failed to unmarshal response body: invalid character 'H' looking for beginning of value`)
		s.Equal(false, resp.Ok)
	})
}
