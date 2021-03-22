package handlers

import (
	"bytes"
	"errors"
	"github.com/odpf/siren/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AlertmanagerServiceMock struct {
	mock.Mock
}

func (m AlertmanagerServiceMock) Migrate() error {
	args := m.Called()
	return args.Error(0)
}

func (m AlertmanagerServiceMock) Upsert(credential domain.AlertCredential) error {
	args := m.Called(credential)

	return args.Error(0)
}

func (m AlertmanagerServiceMock) Get(teamName string) (domain.AlertCredential, error) {
	args := m.Called(teamName)
	return args.Get(0).(domain.AlertCredential), args.Error(1)
}

func TestUpdateAlertCredentials(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("Upsert", mock.Anything).Return(nil)

		payload := []byte(`{
  "entity": "string",
  "pagerduty_credentials": "string",
  "slack_config": {
    "critical": {
      "channel": "string",
      "username": "string",
      "webhook": "http://critical.com"
    },
    "warning": {
      "channel": "string",
      "username": "string",
      "webhook": "http://warning.com"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/alertingCredentials/teams/myTeam", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		updateCredsFn := UpdateAlertCredentials(alertmanagerServiceMock)
		updateCredsFn.ServeHTTP(w, r)
		assert.Equal(t, 201, w.Code)
	})
	t.Run("should return 500 on error", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("Upsert", mock.Anything).Return(errors.New("error occurred while updating"))

		payload := []byte(`{
  "entity": "string",
  "pagerduty_credentials": "string",
  "slack_config": {
    "critical": {
      "channel": "string",
      "username": "string",
      "webhook": "http://critical.com"
    },
    "warning": {
      "channel": "string",
      "username": "string",
      "webhook": "http://warning.com"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/alertingCredentials/teams/myTeam", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		updateCredsFn := UpdateAlertCredentials(alertmanagerServiceMock)
		updateCredsFn.ServeHTTP(w, r)
		assert.Equal(t, 500, w.Code)
	})
	t.Run("should return 4xx on bad  webhook", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("UpsertSlack", mock.Anything).Return(nil)

		payload := []byte(`{
  "entity": "string",
  "pagerduty_credentials": "string",
  "slack_config": {
    "critical": {
      "channel": "string",
      "username": "string",
      "webhook": ":critical"
    },
    "warning": {
      "channel": "string",
      "username": "string",
      "webhook": "http://warning.com"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/alertingCredentials/teams/myTeam", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		updateCredsFn := UpdateAlertCredentials(alertmanagerServiceMock)
		updateCredsFn.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "slack critical webhook is not a valid url", w.Body.String())

	})
	t.Run("should return 4xx on empty entity", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("UpsertSlack", mock.Anything).Return(nil)

		payload := []byte(`{
  "entity": "",
  "pagerduty_credentials": "string",
  "slack_config": {
    "critical": {
      "channel": "string",
      "username": "string",
      "webhook": "critical.com"
    },
    "warning": {
      "channel": "string",
      "username": "string",
      "webhook": "http://warning.com"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/alertingCredentials/teams/myTeam", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		updateCredsFn := UpdateAlertCredentials(alertmanagerServiceMock)
		updateCredsFn.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "entity cannot be empty", w.Body.String())

	})

	t.Run("should return 4xx on empty pagerduty key", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("UpsertSlack", mock.Anything).Return(nil)

		payload := []byte(`{
  "entity": "ssd",
  "pagerduty_credentials": "",
  "slack_config": {
    "critical": {
      "channel": "string",
      "username": "string",
      "webhook": "critical.com"
    },
    "warning": {
      "channel": "string",
      "username": "string",
      "webhook": "http://warning.com"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/alertingCredentials/teams/myTeam", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		updateCredsFn := UpdateAlertCredentials(alertmanagerServiceMock)
		updateCredsFn.ServeHTTP(w, r)
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, "pagerduty key cannot be empty", w.Body.String())

	})

}
