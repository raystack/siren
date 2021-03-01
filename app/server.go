package app

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"log"
	"net/http"

	"github.com/odpf/siren/api"
	"github.com/odpf/siren/store"
)

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	db, err := store.New(&c.DB)
	if err != nil {
		return err
	}

	models := []interface{}{}
	store.Migrate(db, models...)

	r := api.New()

	log.Printf("running server on port %d\n", c.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
}
