package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odpf/siren/api"
	"github.com/odpf/siren/store"
)

// RunServer runs the application server
func RunServer(c *Config) error {
	db, err := store.New(&store.Config{
		Host:     c.DB.Host,
		User:     c.DB.User,
		Password: c.DB.Password,
		Name:     c.DB.Name,
		Port:     c.DB.Port,
		SslMode:  c.DB.SslMode,
	})
	if err != nil {
		return err
	}

	models := []interface{}{}
	store.Migrate(db, models...)

	r := api.New()

	log.Printf("running server on port %d\n", c.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
}
