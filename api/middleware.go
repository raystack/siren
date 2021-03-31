package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func applyMiddlewaresToHandler(middlewares []mux.MiddlewareFunc, next http.Handler) http.Handler {
	handler := next
	for index := range middlewares {
		handler = middlewares[len(middlewares)-1-index](handler)
	}
	return handler
}
