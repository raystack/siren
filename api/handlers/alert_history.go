package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
)

// CreateAlertHistory handler
func CreateAlertHistory(service domain.AlertHistoryService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var alerts domain.Alerts
		err := json.NewDecoder(r.Body).Decode(&alerts)
		if err != nil {
			badRequest(w, err, logger)
			return
		}
		result, err := service.Create(&alerts)

		if err != nil && err.Error() == "alert history parameters missing" {
			badRequest(w, err, logger)
			return
		}
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, result)
		return
	}
}

// GetAlertHistory handler
func GetAlertHistory(service domain.AlertHistoryService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("resource")
		startTimeString := r.URL.Query().Get("startTime")
		endTimeString := r.URL.Query().Get("endTime")
		startTime, err := strconv.ParseUint(startTimeString, 10, 32)
		endTime, err := strconv.ParseUint(endTimeString, 10, 32)
		if name == "" {
			badRequest(w, errors.New("resource query param is required"), logger)
			return
		}
		alerts, err := service.Get(name, uint32(startTime), uint32(endTime))
		if err != nil {
			internalServerError(w, err, logger)
			return
		}
		returnJSON(w, alerts)
		return
	}
}
