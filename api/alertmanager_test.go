package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/logger"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
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

func getPanicLogger() *zap.Logger {
	panicLogger, _ := logger.New(&domain.LogConfig{Level: "panic"})
	return panicLogger
}

func TestGetAlertCredentials(t *testing.T) {
	t.Run("should return alert credentials of the team", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel:  "critical",
					Webhook:  "http://critical.com",
					Username: "critical_user",
				},
				Warning: domain.SlackCredential{
					Channel:  "warning",
					Webhook:  "http://warning.com",
					Username: "warning_user",
				},
			},
		}
		alertmanagerServiceMock.On("Get", "hydra").Return(credential, nil)
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		r, err := http.NewRequest(http.MethodGet, "/alertingCredentials/teams/hydra", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		expectedJson := `{
    "entity": "avengers",
    "team_name": "hydra",
    "pagerduty_credentials": "xyz",
    "slack_config": {
        "critical": {
            "webhook": "http://critical.com",
            "channel": "critical",
            "username": "critical_user"
        },
        "warning": {
            "webhook": "http://warning.com",
            "channel": "warning",
            "username": "warning_user"
        }
    }
}`
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, expectedJson, w.Body.String())
	})

	t.Run("get alert credentials should return 500 on error", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("Get", "hydra").Return(
			domain.AlertCredential{}, errors.New("internal error"))
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		r, err := http.NewRequest(http.MethodGet, "/alertingCredentials/teams/hydra", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		assert.Equal(t, 500, w.Code)
		assert.JSONEq(t, `{"code":500,"message":"Internal server error","data":null}`, w.Body.String())
	})
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
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
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
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
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
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		expectedBody := `{"code":400,"message":"Key: 'AlertCredential.SlackConfig.Critical.Webhook' Error:Field validation for 'Webhook' failed on the 'webhookChecker' tag","data":null}`
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, expectedBody, w.Body.String())

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
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		expectedBody := `{"code":400,"message":"Key: 'AlertCredential.Entity' Error:Field validation for 'Entity' failed on the 'required' tag","data":null}`
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, expectedBody, w.Body.String())

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
		router := New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		expectedBody := `{"code":400,"message":"Key: 'AlertCredential.PagerdutyCredentials' Error:Field validation for 'PagerdutyCredentials' failed on the 'required' tag","data":null}`
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, expectedBody, w.Body.String())

	})
}
