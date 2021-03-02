package handlers

import (
	"encoding/json"
	"github.com/odpf/siren/service"
	"github.com/odpf/siren/store"
	"net/http"
)

// UpsertTemplates handler
func UpsertTemplates(service *service.Container) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var template store.Template
		err := json.NewDecoder(r.Body).Decode(&template)
		if err != nil {
			BadRequest(w, err)
			return
		}
		createdTemplate, err := service.TemplatesService.Upsert(&template)
		if err != nil {
			InternalServerError(w, err)
			return
		}
		returnJSON(w, createdTemplate)
		return
	}
}
