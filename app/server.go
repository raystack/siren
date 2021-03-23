package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odpf/siren/api"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/metric"
	"github.com/odpf/siren/service"
	"github.com/odpf/siren/store"
)

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	nr, err := metric.New(&c.NewRelic)
	if err != nil {
		return err
	}

	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}
	services, err := service.Init(store, c.Cortex, c.Alertmanager)
	if err != nil {
		return err
	}
	r := api.New(services, nr)

	log.Printf("running server on port %d\n", c.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
}

func RunMigrations(c *domain.Config) error {
	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}
	service.MigrateAll(store, *c)
	return nil
}
