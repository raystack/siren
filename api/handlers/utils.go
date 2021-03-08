package handlers

import (
	"encoding/json"
	"net/http"
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

func internalServerError(w http.ResponseWriter, err error) (responseError, error) {
	return sendError(w, defaultErrorMessage, http.StatusInternalServerError, nil)
}

func badRequest(w http.ResponseWriter, err error) {
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	} else {
		errMessage = badRequestErrorMessage
	}

	sendError(w, errMessage, http.StatusBadRequest, nil)
}

func NotFound(w http.ResponseWriter, err error) {
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	} else {
		errMessage = notFoundErrorMessage
	}

	sendError(w, errMessage, http.StatusNotFound, nil)
}

func sendError(w http.ResponseWriter, errorMessage string, code int, data interface{}) (responseError, error) {
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
		return response, err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(jsonBytes)
	if err != nil {
		return response, err
	}

	return response, nil
}
