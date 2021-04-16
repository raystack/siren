package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAlertHistory_CreateAlertHistory(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}

		domainAlerts := domain.Alerts{Alerts: []domain.Alert{{Status: "firing",
			Labels:      domain.Labels{Severity: "CRITICAL"},
			Annotations: domain.Annotations{Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "30"}},
		}}

		payload := []byte(`{"alerts":[{"status":"firing","labels":{"severity":"CRITICAL"},
					"annotations":{"resource":"foo","template":"bar","metricName":"baz","metricValue":"30"}}]}`)

		mockedAlertHistoryService.On("Create", &domainAlerts).Return(dummyAlerts, nil).Once()
		r, err := http.NewRequest(http.MethodPost, "/alertHistory", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.CreateAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyAlerts)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 on bad request", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}

		domainAlerts := domain.Alerts{Alerts: []domain.Alert{{Status: "firing",
			Labels:      domain.Labels{Severity: "CRITICAL"},
			Annotations: domain.Annotations{Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "30"}},
		}}

		payload := []byte(`bad input`)

		mockedAlertHistoryService.On("Create", &domainAlerts).Return(dummyAlerts, nil).Once()
		r, err := http.NewRequest(http.MethodPost, "/alertHistory", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.CreateAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"invalid character 'b' looking for beginning of value\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedAlertHistoryService.AssertNotCalled(t, "Create")
	})

	t.Run("should return 500 on error from service", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		domainAlerts := domain.Alerts{Alerts: []domain.Alert{{Status: "firing",
			Labels:      domain.Labels{Severity: "CRITICAL"},
			Annotations: domain.Annotations{Resource: "foo", Template: "bar", MetricName: "baz", MetricValue: "30"}},
		}}

		payload := []byte(`{"alerts":[{"status":"firing","labels":{"severity":"CRITICAL"},
					"annotations":{"resource":"foo","template":"bar","metricName":"baz","metricValue":"30"}}]}`)

		mockedAlertHistoryService.On("Create", &domainAlerts).
			Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodPost, "/alertHistory", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.CreateAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}

func TestAlertHistory_GetAlertHistory(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}

		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).Return(dummyAlerts, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/alertHistory", nil)
		if err != nil {
			t.Fatal(err)
		}
		q := r.URL.Query()
		q.Add("resource", "foo")
		q.Add("startTime", "100")
		q.Add("endTime", "200")
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		handler := handlers.GetAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyAlerts)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedAlertHistoryService.AssertCalled(t, "Get", "foo", uint32(100), uint32(200))
	})

	t.Run("should return 400 if resource query param is missing", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		dummyAlerts := []domain.AlertHistoryObject{{
			ID: 1, Name: "foo", TemplateID: "bar", MetricName: "bar", MetricValue: "30", Level: "CRITICAL",
		}}

		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).Return(dummyAlerts, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/alertHistory", nil)
		if err != nil {
			t.Fatal(err)
		}
		q := r.URL.Query()
		q.Add("startTime", "100")
		q.Add("endTime", "200")
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		handler := handlers.GetAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"resource query param is required\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedAlertHistoryService.AssertNotCalled(t, "Get")
	})

	t.Run("should return 500 on error from service", func(t *testing.T) {
		mockedAlertHistoryService := &mocks.AlertHistoryService{}
		mockedAlertHistoryService.On("Get", "foo", uint32(100), uint32(200)).
			Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodGet, "/alertHistory", nil)
		if err != nil {
			t.Fatal(err)
		}
		q := r.URL.Query()
		q.Add("resource", "foo")
		q.Add("startTime", "100")
		q.Add("endTime", "200")
		r.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()
		handler := handlers.GetAlertHistory(mockedAlertHistoryService, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedAlertHistoryService.AssertCalled(t, "Get", "foo", uint32(100), uint32(200))
	})
}
