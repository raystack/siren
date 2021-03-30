package handlers

import (
	"embed"
	"net/http"
)

//go:embed swagger.yaml
var swaggerFile embed.FS

func SwaggerFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(swaggerFile)).ServeHTTP(w, r)
		return
	}
}
