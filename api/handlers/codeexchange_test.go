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

func TestCodeExchange_ExchangeCode(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		dummyPayload := domain.OAuthPayload{Code: "foo", Workspace: "bar"}
		dummyResult := domain.OAuthExchangeResponse{OK: true}
		payload := []byte(`{"code": "foo","workspace": "bar"}`)

		mockedCodeExchangeService.On("Exchange", dummyPayload).Return(&dummyResult, nil).Once()
		r, err := http.NewRequest(http.MethodPost, "/code_exchange", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.ExchangeCode(mockedCodeExchangeService, getPanicLogger())
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyResult)
		expectedStringBody := string(response) + "\n"
		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 500 on service error", func(t *testing.T) {
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		dummyPayload := domain.OAuthPayload{Code: "foo", Workspace: "bar"}
		payload := []byte(`{"code": "foo","workspace": "bar"}`)

		mockedCodeExchangeService.On("Exchange", dummyPayload).
			Return(nil, errors.New("random error")).Once()
		r, err := http.NewRequest(http.MethodPost, "/code_exchange", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.ExchangeCode(mockedCodeExchangeService, getPanicLogger())
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 on bad request error", func(t *testing.T) {
		mockedCodeExchangeService := &mocks.CodeExchangeService{}
		payload := []byte(`foobar`)
		r, err := http.NewRequest(http.MethodPost, "/code_exchange", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.ExchangeCode(mockedCodeExchangeService, getPanicLogger())
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"invalid character 'o' in literal false (expecting 'a')\",\"data\":null}"
		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}
