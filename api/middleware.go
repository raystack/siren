package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func logger(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

func SwaggerMiddleware(next http.Handler) http.Handler {
	return middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/swagger.yaml",
		Path: "docs",
	}, next)
}