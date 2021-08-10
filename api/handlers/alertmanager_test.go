package handlers_test

import (
	"bytes"
	"errors"
	"github.com/odpf/siren/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetAlertCredentials(t *testing.T) {
	t.Run("should return alert credentials of the team", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		credential := domain.AlertCredential{
			Entity:               "avengers",
			TeamName:             "hydra",
			PagerdutyCredentials: "xyz",
			SlackConfig: domain.SlackConfig{
				Critical: domain.SlackCredential{
					Channel: "critical",
				},
				Warning: domain.SlackCredential{
					Channel: "warning",
				},
			},
		}
		alertmanagerServiceMock.On("Get", "hydra").Return(credential, nil)
		router := api.New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		r, err := http.NewRequest(http.MethodGet, "/teams/hydra/credentials", nil)
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
            "channel": "critical"
        },
        "warning": {
            "channel": "warning"
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
		router := api.New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		r, err := http.NewRequest(http.MethodGet, "/teams/hydra/credentials", nil)
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
      "channel": "string"
    },
    "warning": {
      "channel": "string"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/teams/myTeam/credentials", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router := api.New(&service.Container{
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
      "channel": "string"
    },
    "warning": {
      "channel": "string"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/teams/myTeam/credentials", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router := api.New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		assert.Equal(t, 500, w.Code)
	})

	t.Run("should return 400 on invalid json payload", func(t *testing.T) {
		var alertmanagerServiceMock AlertmanagerServiceMock
		alertmanagerServiceMock.On("UpsertSlack", mock.Anything).Return(nil)

		payload := []byte(`abcd`)
		r, err := http.NewRequest(http.MethodPut, "/teams/myTeam/credentials", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router := api.New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		expectedBody := `{"code":400,"message":"invalid character 'a' looking for beginning of value","data":null}`
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
      "channel": "string"
    },
    "warning": {
      "channel": "string"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/teams/myTeam/credentials", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router := api.New(&service.Container{
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
      "channel": "string"
    },
    "warning": {
      "channel": "string"
    }
  }}`)
		r, err := http.NewRequest(http.MethodPut, "/teams/myTeam/credentials", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router := api.New(&service.Container{
			AlertmanagerService: alertmanagerServiceMock,
		}, nil, getPanicLogger())
		router.ServeHTTP(w, r)
		expectedBody := `{"code":400,"message":"Key: 'AlertCredential.PagerdutyCredentials' Error:Field validation for 'PagerdutyCredentials' failed on the 'required' tag","data":null}`
		assert.Equal(t, 400, w.Code)
		assert.Equal(t, expectedBody, w.Body.String())

	})
}
