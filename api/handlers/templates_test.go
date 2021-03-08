package handlers_test

import (
	"bytes"
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

func TestTemplates_UpsertTemplates(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		}

		payload := []byte(`{"name":"foo", "body": "bar", "tags": ["test"], "variables": [{"name": "test-name", "default":"test-default", "description": "test-description", "type": "test-type" }]}`)

		mockedTemplatesService.On("Upsert", dummyTemplate).Return(dummyTemplate, nil).Once()
		r, err := http.NewRequest(http.MethodPut, "/templates", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyTemplate)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Upsert", dummyTemplate)
	})

	t.Run("should return 400 Bad Request on failure", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		payload := []byte(`{"foo"}`)
		r, err := http.NewRequest(http.MethodPut, "/templates", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"invalid character '}' after object key\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 400 Bad Request if name is empty in request body", func(t *testing.T) {
		expectedError := errors.New("name cannot be empty")
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplate := &domain.Template{
			Name: "", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		}

		payload := []byte(`{"name":"", "body": "bar", "tags": ["test"], "variables": [{"name": "test-name", "default":"test-default", "description": "test-description", "type": "test-type" }]}`)

		mockedTemplatesService.On("Upsert", dummyTemplate).Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodPut, "/templates", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"name cannot be empty\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Upsert", dummyTemplate)
	})

	t.Run("should return 400 Bad Request if body is empty in request body", func(t *testing.T) {
		expectedError := errors.New("body cannot be empty")
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		}

		payload := []byte(`{"name":"foo", "body": "", "tags": ["test"], "variables": [{"name": "test-name", "default":"test-default", "description": "test-description", "type": "test-type" }]}`)

		mockedTemplatesService.On("Upsert", dummyTemplate).Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodPut, "/templates", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"body cannot be empty\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Upsert", dummyTemplate)
	})

	t.Run("should return 500 Error on failure", func(t *testing.T) {
		expectedError := errors.New("random error")
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplate := &domain.Template{
			Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		}

		payload := []byte(`{"name":"foo", "body": "bar", "tags": ["test"], "variables": [{"name": "test-name", "default":"test-default", "description": "test-description", "type": "test-type" }]}`)

		mockedTemplatesService.On("Upsert", dummyTemplate).Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodPut, "/templates", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.UpsertTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Upsert", dummyTemplate)
	})
}

func TestTemplates_GetTemplates(t *testing.T) {
	t.Run("should return 200 OK on success if template exist", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplate := &domain.Template{
			ID: 1, Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		}
		mockedTemplatesService.On("GetByName", "foo").Return(dummyTemplate, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyTemplate)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "GetByName", "foo")
	})

	t.Run("should return 404 Not found if template not exist", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		mockedTemplatesService.On("GetByName", "foo").Return(nil, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusNotFound
		expectedStringBody := "{\"code\":404,\"message\":\"not found\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 500 Error on any failure", func(t *testing.T) {
		expectedError := errors.New("random error")
		mockedTemplatesService := &mocks.TemplatesService{}
		mockedTemplatesService.On("GetByName", "foo").Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.GetTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}

func TestTemplates_IndexTemplates(t *testing.T) {
	t.Run("should return 200 OK on success for non-empty tag", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		dummyTemplates := []domain.Template{{
			ID: 1, Name: "foo", Body: "bar",
			Tags: []string{"test"},
			Variables: []domain.Variable{{
				Name:        "test-name",
				Default:     "test-default",
				Description: "test-description",
				Type:        "test-type",
			}},
		},
		}
		mockedTemplatesService.On("Index", "foo").Return(dummyTemplates, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		q := r.URL.Query()
		q.Add("tag", "foo")
		r.URL.RawQuery = q.Encode()
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.IndexTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyTemplates)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Index", "foo")
	})

	t.Run("should return 200 OK on success for empty tag", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		var dummyTemplates []domain.Template
		mockedTemplatesService.On("Index", "").Return(dummyTemplates, nil).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.IndexTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		response, _ := json.Marshal(dummyTemplates)
		expectedStringBody := string(response) + "\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Index", "")
	})

	t.Run("should return 500 Error on any failure", func(t *testing.T) {
		expectedError := errors.New("random error")
		mockedTemplatesService := &mocks.TemplatesService{}
		mockedTemplatesService.On("Index", "").Return(nil, expectedError).Once()
		r, err := http.NewRequest(http.MethodGet, "/templates", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.IndexTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}

func TestTemplates_RenderTemplates(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		inputBody := make(map[string]string)
		inputBody["foo"] = "bar"
		payload := []byte(`{"foo":"bar"}`)
		mockedTemplatesService.On("Render", "foo", inputBody).Return("foo bar baz", nil).Once()
		r, err := http.NewRequest(http.MethodPost, "/templates/{name}/render", bytes.NewBuffer(payload))
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.RenderTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		expectedStringBody := "\"foo bar baz\"\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Render", "foo", inputBody)
	})

	t.Run("should return 400 Bad Request if bad input given", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		badPayload := []byte(`{"foo"}`)
		r, err := http.NewRequest(http.MethodPost, "/templates/{name}/render", bytes.NewBuffer(badPayload))
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.RenderTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusBadRequest
		expectedStringBody := "{\"code\":400,\"message\":\"invalid character '}' after object key\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 404 Not found if template not exist", func(t *testing.T) {
		expectedError := errors.New("template not found")
		mockedTemplatesService := &mocks.TemplatesService{}
		inputBody := make(map[string]string)
		inputBody["foo"] = "bar"
		payload := []byte(`{"foo":"bar"}`)
		mockedTemplatesService.On("Render", "foo", inputBody).Return("", expectedError).Once()
		r, err := http.NewRequest(http.MethodPost, "/templates/{name}/render", bytes.NewBuffer(payload))
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.RenderTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusNotFound
		expectedStringBody := "{\"code\":404,\"message\":\"template not found\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})

	t.Run("should return 500 Error on any failure", func(t *testing.T) {
		expectedError := errors.New("random error")
		mockedTemplatesService := &mocks.TemplatesService{}
		inputBody := make(map[string]string)
		inputBody["foo"] = "bar"
		payload := []byte(`{"foo":"bar"}`)
		mockedTemplatesService.On("Render", "foo", inputBody).Return("", expectedError).Once()
		r, err := http.NewRequest(http.MethodPost, "/templates/{name}/render", bytes.NewBuffer(payload))
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.RenderTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}

func TestTemplates_DeleteTemplates(t *testing.T) {
	t.Run("should return 200 OK on success", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		mockedTemplatesService.On("Delete", "foo").Return(nil).Once()
		r, err := http.NewRequest(http.MethodDelete, "/templates", nil)
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.DeleteTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusOK
		expectedStringBody := "null\n"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
		mockedTemplatesService.AssertCalled(t, "Delete", "foo")
	})

	t.Run("should return 500 Error on any failure", func(t *testing.T) {
		mockedTemplatesService := &mocks.TemplatesService{}
		expectedError := errors.New("random error")
		mockedTemplatesService.On("Delete", "foo").Return(expectedError).Once()
		r, err := http.NewRequest(http.MethodDelete, "/templates", nil)
		r = mux.SetURLVars(r, map[string]string{"name": "foo"})
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()
		handler := handlers.DeleteTemplates(mockedTemplatesService)
		expectedStatusCode := http.StatusInternalServerError
		expectedStringBody := "{\"code\":500,\"message\":\"Internal server error\",\"data\":null}"

		handler.ServeHTTP(w, r)

		assert.Equal(t, expectedStatusCode, w.Code)
		assert.Equal(t, expectedStringBody, w.Body.String())
	})
}
