package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
)

// ExchangeCode handler
func ExchangeCode(service domain.CodeExchangeService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload domain.OAuthPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			badRequest(w, err, logger)
			return
		}
		result, err := service.Exchange(payload)

		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, result)
		return
	}
}
