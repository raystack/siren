package handlers

import (
	"encoding/json"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/slacknotifier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
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

			err = payload.Validate()
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
			result, err := notifierServices.Slack.Notify(&payload)
			if err != nil {
				if isBadRequest(err) {
					badRequest(w, err, logger)
					return
				}
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

func isBadRequest(err error) bool {

	var noChannelFoundError *slacknotifier.NoChannelFoundErr
	var userLookupByEmailErr *slacknotifier.UserLookupByEmailErr

	return errors.As(err, &noChannelFoundError) || errors.As(err, &userLookupByEmailErr)
}
