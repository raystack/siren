package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/odpf/siren/api/handlers"
	"github.com/odpf/siren/service"
)

// New initializes the service router
func New(container *service.Container) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Use(logger)

	// Route => handler
	r.Methods("GET").Path("/swagger.yaml").Handler(handlers.SwaggerFile())
	r.Methods("GET").Path("/docs").Handler(middleware.SwaggerUI(middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}, r.NotFoundHandler))

	r.Methods("GET").Path("/ping").Handler(handlers.Ping())

	r.Methods("PUT").Path("/templates").Handler(handlers.UpsertTemplates(container.TemplatesService))
	r.Methods("GET").Path("/templates").Handler(handlers.IndexTemplates(container.TemplatesService))
	r.Methods("GET").Path("/templates/{name}").Handler(handlers.GetTemplates(container.TemplatesService))
	r.Methods("DELETE").Path("/templates/{name}").Handler(handlers.DeleteTemplates(container.TemplatesService))
	r.Methods("POST").Path("/templates/{name}/render").Handler(handlers.RenderTemplates(container.TemplatesService))

	return r
}
