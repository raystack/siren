package httpreceiver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/siren/pkg/retry"
	"github.com/goto/siren/plugins/receivers/httpreceiver"
	"github.com/stretchr/testify/assert"
)

func TestService_Notify_WithoutRetrier(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     context.Context
		cfg     httpreceiver.AppConfig
		apiURL  string
		message []byte
		wantErr bool
	}{
		{
			name:    "should return error if json marshal error",
			message: []byte(`{//`),
			wantErr: true,
		},
		{
			name:    "should return error if create request return error",
			apiURL:  "http://localhost",
			wantErr: true,
		},
		{
			name:    "should return error if there is failure in http call",
			ctx:     context.Background(),
			apiURL:  "xxx",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := httpreceiver.NewPluginService(log.NewNoop(), tc.cfg)
			if err := c.Notify(tc.ctx, tc.apiURL, tc.message); (err != nil) != tc.wantErr {
				t.Errorf("Client.Notify() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestService_Notify_HTTPCall(t *testing.T) {
	t.Run("should return error if error response is retryable", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTooManyRequests)
		}))

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{})
		err := c.Notify(context.Background(), testServer.URL, nil)

		assert.EqualError(t, err, "Too Many Requests")

		testServer.Close()
	})

	t.Run("should return error if error response is non retryable", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{})
		err := c.Notify(context.Background(), testServer.URL, nil)

		assert.EqualError(t, err, "Bad Request")

		testServer.Close()
	})

	t.Run("should return error if read response body error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{})
		err := c.Notify(context.Background(), testServer.URL, nil)

		assert.EqualError(t, err, "failed to read response body: unexpected EOF")

		testServer.Close()
	})

	t.Run("should return error if read response body error", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "1")
		}))

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{})
		err := c.Notify(context.Background(), testServer.URL, nil)

		assert.EqualError(t, err, "failed to read response body: unexpected EOF")

		testServer.Close()
	})
}

func TestService_Notify_WithRetrier(t *testing.T) {
	var expectedCounter = 4

	t.Run("when 429 is returned", func(t *testing.T) {
		var counter = 0
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter = counter + 1
			w.Header().Set("Retry-After", "10")
			w.WriteHeader(http.StatusTooManyRequests)
		}))

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{},
			httpreceiver.WithRetrier(retry.New(retry.Config{Enable: true})))
		_ = c.Notify(context.Background(), testServer.URL, nil)

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

		c := httpreceiver.NewPluginService(log.NewNoop(), httpreceiver.AppConfig{},
			httpreceiver.WithRetrier(retry.New(retry.Config{Enable: true})))
		_ = c.Notify(context.Background(), testServer.URL, nil)

		assert.Equal(t, expectedCounter, counter)

		testServer.Close()
	})

}
