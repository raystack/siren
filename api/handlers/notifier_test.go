package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotifier_Notify(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		dummyMessage := domain.SlackMessage{ReceiverName: "foo",
			ReceiverType: "user",
			Message:      "some text",
			Entity:       "odpf"}

		payload := []byte(`{"receiver_name": "foo","receiver_type": "user","entity": "odpf","message": "some text"}`)
		expectedResponse := domain.SlackMessageSendResponse{
			OK: true,
		}
		mockedSlackNotifierService.On("Notify", &dummyMessage).Return(&expectedResponse, nil).Once()
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=slack", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(expectedResponse)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 500 Internal server error on failure", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		dummyMessage := domain.SlackMessage{ReceiverName: "foo",
			ReceiverType: "user",
			Message:      "some text",
			Entity:       "odpf"}

		payload := []byte(`{"receiver_name": "foo","receiver_type": "user","entity": "odpf","message": "some text"}`)
		mockedSlackNotifierService.On("Notify", &dummyMessage).Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=slack", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := `{"code":500,"message":"Internal server error","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad request if app not part of channel", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		dummyMessage := domain.SlackMessage{ReceiverName: "test",
			ReceiverType: "channel",
			Message:      "some text",
			Entity:       "odpf"}
		expectedError := errors.New("app is not part of test")
		payload := []byte(`{"receiver_name": "test","receiver_type": "channel","entity": "odpf","message": "some text"}`)
		mockedSlackNotifierService.On("Notify", &dummyMessage).Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=slack", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := `{"code":400,"message":"app is not part of test","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad request if user not found", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		dummyMessage := domain.SlackMessage{ReceiverName: "foo",
			ReceiverType: "user",
			Message:      "some text",
			Entity:       "odpf"}
		expectedError := errors.New("failed to get id for foo")
		payload := []byte(`{"receiver_name": "foo","receiver_type": "user","entity": "odpf","message": "some text"}`)
		mockedSlackNotifierService.On("Notify", &dummyMessage).Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=slack", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := `{"code":400,"message":"failed to get id for foo","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad request if no provider specified", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		payload := []byte(`{"receiver_name": "foo","receiver_type": "user","entity": "odpf","message": "some text"}`)
		r, err := http.NewRequest(http.MethodPost, "/notifications", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := `{"code":400,"message":"provider not given in query params","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad request if unknown provider specified", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		payload := []byte(`{"receiver_name": "foo","receiver_type": "user","entity": "odpf","message": "some text"}`)
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=email", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := `{"code":400,"message":"unrecognized provider","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad request for invalid payload", func(t *testing.T) {
		mockedSlackNotifierService := &mocks.SlackNotifierService{}
		notifierServices := domain.NotifierServices{Slack: mockedSlackNotifierService}
		payload := []byte(`abcd`)
		r, err := http.NewRequest(http.MethodPost, "/notifications?provider=slack", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.Notify(notifierServices, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := `{"code":400,"message":"invalid character 'a' looking for beginning of value","data":null}`
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}
