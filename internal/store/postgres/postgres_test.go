package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/gorm"
)

const logLevelDebug = "debug"

var (
	pgConfig = postgres.Config{
		Host:     "localhost",
		User:     "test_user",
		Password: "test_pass",
		Name:     "test_db",
		SSLMode:  "disable",
		LogLevel: "debug",
	}
)

func newTestClient(logger log.Logger) (*gorm.DB, *dockertest.Pool, *dockertest.Resource, error) {

	opts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_PASSWORD=" + pgConfig.Password,
			"POSTGRES_USER=" + pgConfig.User,
			"POSTGRES_DB=" + pgConfig.Name,
		},
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not create dockertest pool: %w", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(opts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	pgConfig.Port = resource.GetPort("5432/tcp")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot parse external port of container to int: %w", err)
	}

	// attach terminal logger to container if exists
	// for debugging purpose
	if logger.Level() == logLevelDebug {
		logWaiter, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
			Container:    resource.Container.ID,
			OutputStream: logger.Writer(),
			ErrorStream:  logger.Writer(),
			Stderr:       true,
			Stdout:       true,
			Stream:       true,
		})
		if err != nil {
			logger.Fatal("could not connect to postgres container log output", "error", err)
		}
		defer func() {
			err = logWaiter.Close()
			if err != nil {
				logger.Fatal("could not close container log", "error", err)
			}

			err = logWaiter.Wait()
			if err != nil {
				logger.Fatal("could not wait for container log to close", "error", err)
			}
		}()
	}

	// Tell docker to hard kill the container in 120 seconds
	if err := resource.Expire(120); err != nil {
		return nil, nil, nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 60 * time.Second

	var pgClient *gorm.DB
	if err = pool.Retry(func() error {
		pgClient, err = postgres.New(logger, pgConfig)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	err = setup(context.Background(), logger, pgClient)
	if err != nil {
		logger.Fatal("failed to setup and migrate DB", "error", err)
	}
	return pgClient, pool, resource, nil
}

func purgeDocker(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %w", err)
	}
	return nil
}

func setup(ctx context.Context, logger log.Logger, db *gorm.DB) (err error) {
	var queries = []string{
		"DROP SCHEMA public CASCADE",
		"CREATE SCHEMA public",
	}

	err = execQueries(ctx, db, queries)
	if err != nil {
		return
	}

	err = postgres.Migrate(logger, db)
	return
}

// ExecQueries is used for executing list of db query
func execQueries(ctx context.Context, db *gorm.DB, queries []string) error {
	for _, query := range queries {
		err := db.Exec(query).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func bootstrapProvider(db *gorm.DB) ([]provider.Provider, error) {
	filePath := "./testdata/mock-provider.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []provider.Provider
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []provider.Provider
	for _, d := range data {
		var providerModel model.Provider
		if err := providerModel.FromDomain(&d); err != nil {
			return nil, err
		}

		result := db.Create(&providerModel)
		if result.Error != nil {
			return nil, result.Error
		}

		insertedProvider, err := providerModel.ToDomain()
		if err != nil {
			return nil, err
		}
		insertedData = append(insertedData, *insertedProvider)
	}

	return insertedData, nil
}
