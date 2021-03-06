package api

import (
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func logger(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
