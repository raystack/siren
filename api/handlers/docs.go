package handlers

import (
	"net/http"
)

func SwaggerFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger.yaml")
		return
	}
}
