package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const (
	defaultErrorMessage    = "Internal server error"
	badRequestErrorMessage = "Bad Request"
	notFoundErrorMessage   = "Not Found"
)

func returnJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type responseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func internalServerError(w http.ResponseWriter, err error, logger *zap.Logger) {
	sendError(w, defaultErrorMessage, http.StatusInternalServerError, nil, logger)
}

func badRequest(w http.ResponseWriter, err error, logger *zap.Logger) {
	var errMessage string
	if err != nil {
		logger.Error("handler", zap.String("error", err.Error()))
		errMessage = err.Error()
	} else {
		errMessage = badRequestErrorMessage
	}

	sendError(w, errMessage, http.StatusBadRequest, nil, logger)
}

func notFound(w http.ResponseWriter, err error, logger *zap.Logger) {
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	} else {
		errMessage = notFoundErrorMessage
	}

	sendError(w, errMessage, http.StatusNotFound, nil, logger)
}

func sendError(w http.ResponseWriter, errorMessage string, code int, data interface{}, logger *zap.Logger) {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	response := responseError{
		Code:    code,
		Message: errorMessage,
		Data:    data,
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		internalServerError(w, err, logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(jsonBytes)
	if err != nil {
		internalServerError(w, err, logger)
		return
	}
	logger.Error("handler", zap.String("error", errorMessage))
}
