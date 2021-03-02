package api

import (
	"github.com/gorilla/mux"
	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/service"
)

// New initializes the service router
func New(container *service.Container) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(logger)

	r.Methods("GET").Path("/ping").Handler(handlers.Ping())
	r.Methods("PUT").Path("/templates").Handler(handlers.UpsertTemplates(container))

	return r
}
