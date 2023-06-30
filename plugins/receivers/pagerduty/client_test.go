package pagerduty_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raystack/siren/pkg/retry"
	"github.com/raystack/siren/plugins/receivers/pagerduty"
	"github.com/stretchr/testify/assert"
)

func TestClient_NotifyV1_WithoutRetrier(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		cfg     pagerduty.AppConfig
		message pagerduty.MessageV1
		wantErr bool
	}{
		{
			name:    "should return error if json marshal error",
			message: pagerduty.MessageV1{},
			wantErr: true,
		},
		{
			name:    "should return error if create request return error",
			cfg:     pagerduty.AppConfig{APIHost: "http://localhost"},
			wantErr: true,
		},
		{
			name:    "should return error if there is failure in http call",
			ctx:     context.Background(),
			cfg:     pagerduty.AppConfig{APIHost: "xxx"},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := pagerduty.NewClient(tc.cfg)
			if err := c.NotifyV1(tc.ctx, tc.message); (err != nil) != tc.wantErr {
				t.Errorf("Client.Notify() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestClient_NotifyV1_HTTPCall(t *testing.T) {
	t.Run("should return error if error response is retryable", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "Too Many Requests")

		testServer.Close()
	})

	t.Run("should return error if error response is non retryable", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "error with status code Bad Request and body ")

		testServer.Close()
	})

	t.Run("should return error if read response body error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "failed to read response body: unexpected EOF")

		testServer.Close()
	})

	t.Run("should return error if read response body error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "failed to read response body: unexpected EOF")

		testServer.Close()
	})

	t.Run("should return error if response can't be unmarshalled", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{//x`))
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "failed to unmarshal response body: invalid character '/' looking for beginning of object key string")

		testServer.Close()
	})

	t.Run("should return error if response status is not success", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"failed"}`))
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.EqualError(t, err, "something wrong when sending pagerduty event: {failed  }")

		testServer.Close()
	})

	t.Run("should return error if response status is not success", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"success"}`))
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL})
		err := c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.NoError(t, err)

		testServer.Close()
	})
}

func TestClient_NotifyV1_WithRetrier(t *testing.T) {
	var expectedCounter = 4

	t.Run("when 429 is returned", func(t *testing.T) {
		var counter = 0
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter + 1
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL},
			pagerduty.ClientWithRetrier(retry.New(retry.Config{Enable: true})),
		)
		_ = c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.Equal(t, expectedCounter, counter)

		testServer.Close()
	})

	t.Run("when 5xx is returned", func(t *testing.T) {
		var counter = 0
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter + 1
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"ok":false}`))
		}))

		c := pagerduty.NewClient(pagerduty.AppConfig{APIHost: testServer.URL},
			pagerduty.ClientWithRetrier(retry.New(retry.Config{Enable: true})),
		)
		_ = c.NotifyV1(context.Background(), pagerduty.MessageV1{})

		assert.Equal(t, expectedCounter, counter)

		testServer.Close()
	})

}
