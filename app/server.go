package app

import (
	"fmt"
	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	"net/http"

	"github.com/odpf/siren/api"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/logger"
	"github.com/odpf/siren/metric"
	"github.com/odpf/siren/service"
	"github.com/odpf/siren/store"
	"go.uber.org/zap"
)

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	nr, err := metric.New(&c.NewRelic)
	if err != nil {
		return err
	}

	logger, err := logger.New(&c.Log)
	if err != nil {
		return err
	}

	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}
	cortexConfig := cortexClient.Config{
		Address:         c.Cortex.Address,
		UseLegacyRoutes: false,
	}
	client, err := cortexClient.New(cortexConfig)
	if err != nil {
		return nil
	}
	services, err := service.Init(store, c.Cortex, c.SirenService, client)
	if err != nil {
		return err
	}
	r := api.New(services, nr, logger)

	logger.Info("starting server", zap.Int("port", c.Port))
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r)
}

func RunMigrations(c *domain.Config) error {
	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}

	cortexConfig := cortexClient.Config{
		Address:         c.Cortex.Address,
		UseLegacyRoutes: false,
	}
	client, err := cortexClient.New(cortexConfig)
	if err != nil {
		return nil
	}
	services, err := service.Init(store, c.Cortex, c.SirenService, client)
	if err != nil {
		return err
	}

	services.MigrateAll(store)
	return nil
}
