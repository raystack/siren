package handlers

import (
	"encoding/json"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"github.com/pkg/errors"
	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
)

// UpsertRule handler
func UpsertRule(service domain.RuleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rule domain.Rule
		err := json.NewDecoder(r.Body).Decode(&rule)
		if err != nil {
			badRequest(w, err, logger)
			return
		}

		err = rule.Validate()
		if err != nil {
			var e *validator.InvalidValidationError
			if errors.As(err, &e) {
				logger.Error("invalid validation error")
				internalServerError(w, err, logger)
				return
			}
			badRequest(w, err, logger)
			return
		}
		upsertedRule, err := service.Upsert(&rule)
		if err != nil && err.Error() == ("template not found") {
			badRequest(w, err, logger)
			return
		}
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, upsertedRule)
		return
	}
}

// GetRules handler
func GetRules(service domain.RuleService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		namespace := r.URL.Query().Get("namespace")
		entity := r.URL.Query().Get("entity")
		groupName := r.URL.Query().Get("group_name")
		status := r.URL.Query().Get("status")
		template := r.URL.Query().Get("template")
		rules, err := service.Get(namespace, entity, groupName, status, template)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, rules)
		return
	}
}
