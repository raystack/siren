package handlers

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/templates"
	"net/http"
)

// UpsertTemplates handler
func UpsertTemplates(service *templates.Service) http.HandlerFunc {
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
