package handlers

import (
	"encoding/json"
	"net/http"
)

const (
	defaultErrorMessage    = "Internal server error"
	BadRequestErrorMessage = "Bad Request"
)

func returnJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InternalServerError(w http.ResponseWriter, err error) (ResponseError, error) {
	return sendError(w, defaultErrorMessage, http.StatusInternalServerError, nil)
}

func BadRequest(w http.ResponseWriter, err error) {
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	} else {
		errMessage = BadRequestErrorMessage
	}

	sendError(w, errMessage, http.StatusBadRequest, nil)
}

func sendError(w http.ResponseWriter, errorMessage string, code int, data interface{}) (ResponseError, error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	response := ResponseError{
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
