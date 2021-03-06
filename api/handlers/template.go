package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/odpf/siren/domain"
	"net/http"
)

// UpsertTemplates handler
func UpsertTemplates(service domain.TemplatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var template domain.Template
		err := json.NewDecoder(r.Body).Decode(&template)
		if err != nil {
			BadRequest(w, err)
			return
		}
		createdTemplate, err := service.Upsert(&template)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		returnJSON(w, createdTemplate)
		return
	}
}

// IndexTemplates handler
func IndexTemplates(service domain.TemplatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.URL.Query().Get("tag")
		var templates []domain.Template
		var err error
		templates, err = service.Index(tag)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		returnJSON(w, templates)
		return
	}
}

// GetTemplates handler
func GetTemplates(service domain.TemplatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		templates, err := service.GetByName(name)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		if templates == nil {
			NotFound(w, fmt.Errorf("not found"))
			return
		}
		returnJSON(w, templates)
		return
	}
}

// DeleteTemplates handler
func DeleteTemplates(service domain.TemplatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		err := service.Delete(name)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		returnJSON(w, nil)
		return
	}
}

// RenderTemplates handler
func RenderTemplates(service domain.TemplatesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name := params["name"]
		var body map[string]string
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			BadRequest(w, err)
			return
		}
		renderedBody, err := service.Render(name, body)
		if err != nil && err.Error() == "template not found" {
			NotFound(w, err)
			return
		}
		if err != nil {
			InternalServerError(w, err)
			return
		}
		returnJSON(w, renderedBody)
		return
	}
}
