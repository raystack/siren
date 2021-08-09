package handlers_test

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkspace_GetWorkspaceChannels(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedWorkspaceService := &mocks.WorkspaceService{}
		dummyResult := []domain.Channel{
			{Name: "foo"},
			{Name: "bar"},
		}
		mockedWorkspaceService.On("GetChannels", "random").Return(dummyResult, nil).Once()

		r, err := http.NewRequest(http.MethodGet, "/workspaces/{workspaceName}/channels", nil)
		r = mux.SetURLVars(r, map[string]string{"workspaceName": "random"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetWorkspaceChannels(mockedWorkspaceService, getPanicLogger())

		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyResult)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedWorkspaceService.AssertCalled(t, "GetChannels", "random")
	})

	t.Run("should return 500 Error on any failure", func(t *testing.T) {
		mockedWorkspaceService := &mocks.WorkspaceService{}
		expectedError := errors.New("random error")
		mockedWorkspaceService.On("GetChannels", "random").Return(nil, expectedError).Once()

		r, err := http.NewRequest(http.MethodGet, "/workspaces/{workspaceName}/channels", nil)
		r = mux.SetURLVars(r, map[string]string{"workspaceName": "random"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetWorkspaceChannels(mockedWorkspaceService, getPanicLogger())

		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)
		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}
