package handlers

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

// Notify handler
func Notify(notifierServices domain.NotifierServices, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.URL.Query().Get("provider")
		switch provider {
		case "slack":
			var payload domain.SlackMessage
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				badRequest(w, err, logger)
				return
			}
			result, err := notifierServices.Slack.Notify(&payload)
			if err != nil {
				internalServerError(w, err, logger)
				return
			}
			returnJSON(w, result)
			return
		case "":
			badRequest(w, errors.New("provider not given in query params"), logger)
			return
		default:
			badRequest(w, errors.New("unrecognized provider"), logger)
			return
		}

	}
}
