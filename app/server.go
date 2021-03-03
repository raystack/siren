package app

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"log"
	"net/http"

	"github.com/odpf/siren/api"
	"github.com/odpf/siren/store"
)

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}
	services := service.Init(store)

	r := api.New(services)

	log.Printf("running server on port %d\n", c.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
}
