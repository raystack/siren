package handlers

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"net/http"
)

// UpsertRule handler
func UpsertRule(service domain.RuleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rule domain.Rule
		err := json.NewDecoder(r.Body).Decode(&rule)
		if err != nil {
			badRequest(w, err)
			return
		}
		upsertedRule, err := service.Upsert(&rule)
		if err != nil {
			internalServerError(w, err)
			return
		}
		returnJSON(w, upsertedRule)
		return
	}
}
