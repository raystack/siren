package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

// UpsertTemplates handler
func UpsertTemplates(service domain.TemplatesService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var template domain.Template
		err := json.NewDecoder(r.Body).Decode(&template)
		if err != nil {
			badRequest(w, err, logger)
			return
		}
		v := validator.New()
		err = v.Struct(template)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				logger.Error("invalid validation error")
				internalServerError(w, err, logger)
				return
			}
			badRequest(w, err, logger)
			return
		}
		createdTemplate, err := service.Upsert(&template)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, createdTemplate)
		return
	}
}

// IndexTemplates handler
func IndexTemplates(service domain.TemplatesService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		templates, err := service.Index(tag)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, templates)
		return
	}
}

// GetTemplates handler
func GetTemplates(service domain.TemplatesService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		templates, err := service.GetByName(name)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		if templates == nil {
			notFound(w, errors.New(notFoundErrorMessage), logger)
			return
		}
		returnJSON(w, templates)
		return
	}
}

// DeleteTemplates handler
func DeleteTemplates(service domain.TemplatesService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		err := service.Delete(name)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, nil)
		return
	}
}

// RenderTemplates handler
func RenderTemplates(service domain.TemplatesService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		var body map[string]string
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			badRequest(w, err, logger)
			return
		}
		renderedBody, err := service.Render(name, body)
		if err != nil && err.Error() == "template not found" {
			notFound(w, err, logger)
			return
		}
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, renderedBody)
		return
	}
}
