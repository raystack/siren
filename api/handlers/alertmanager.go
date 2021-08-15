package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

func UpdateAlertCredentials(service domain.AlertmanagerService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var alertCredential domain.AlertCredential
		err := json.NewDecoder(r.Body).Decode(&alertCredential)
		if err != nil {
			badRequest(w, err, logger)
			return
		}
		err = alertCredential.Validate()
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
		teamName := params["teamName"]
		alertCredential.TeamName = teamName

		err = service.Upsert(alertCredential)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		w.WriteHeader(201)
		return
	}
}

func GetAlertCredentials(service domain.AlertmanagerService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		teamName := params["teamName"]
		alertCredential, err := service.Get(teamName)
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(alertCredential)
		w.WriteHeader(201)
		return
	}
}
