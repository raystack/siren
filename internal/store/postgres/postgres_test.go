package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.uber.org/zap"
)

const logLevelDebug = "debug"

var (
	pgConfig = postgres.Config{
		Host:     "localhost",
		User:     "test_user",
		Password: "test_pass",
		Name:     "test_db",
		SSLMode:  "disable",
	}
)

func newTestClient(logger log.Logger) (*postgres.Client, *dockertest.Pool, *dockertest.Resource, error) {

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

	var pgClient *postgres.Client
	if err = pool.Retry(func() error {
		pgClient, err = postgres.NewClient(logger, pgConfig)
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

func setup(ctx context.Context, logger log.Logger, client *postgres.Client) (err error) {
	var queries = []string{
		"DROP SCHEMA public CASCADE",
		"CREATE SCHEMA public",
	}

	err = execQueries(ctx, client, queries)
	if err != nil {
		return
	}

	err = client.Migrate()
	return
}

// ExecQueries is used for executing list of db query
func execQueries(ctx context.Context, client *postgres.Client, queries []string) error {
	for _, query := range queries {
		err := client.GetDB(context.TODO()).Exec(query).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func bootstrapProvider(client *postgres.Client) ([]provider.Provider, error) {
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
		var mdl model.Provider
		mdl.FromDomain(&d)

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		insertedData = append(insertedData, *mdl.ToDomain())
	}

	return insertedData, nil
}

func bootstrapNamespace(client *postgres.Client) ([]namespace.EncryptedNamespace, error) {
	filePath := "./testdata/mock-namespace.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []namespace.Namespace
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []namespace.EncryptedNamespace
	for _, d := range data {
		var mdl model.Namespace
		mdl.FromDomain(&namespace.EncryptedNamespace{
			Namespace:   &d,
			Credentials: fmt.Sprintf("%+v", &d.Credentials),
		})

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		insertedData = append(insertedData, *mdl.ToDomain())
	}

	return insertedData, nil
}

func bootstrapReceiver(client *postgres.Client) ([]receiver.Receiver, error) {
	filePath := "./testdata/mock-receiver.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []receiver.Receiver
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []receiver.Receiver
	for _, d := range data {
		var mdl model.Receiver
		mdl.FromDomain(&d)

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		insertedData = append(insertedData, *mdl.ToDomain())
	}

	return insertedData, nil
}

func bootstrapAlert(client *postgres.Client) ([]alert.Alert, error) {
	filePath := "./testdata/mock-alert.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []alert.Alert
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []alert.Alert
	for _, d := range data {
		var mdl model.Alert
		mdl.FromDomain(&d)

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		dmn := mdl.ToDomain()
		insertedData = append(insertedData, *dmn)
	}

	return insertedData, nil
}

func bootstrapTemplate(client *postgres.Client) ([]template.Template, error) {
	filePath := "./testdata/mock-template.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []template.Template
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []template.Template
	for _, d := range data {
		var mdl model.Template
		if err := mdl.FromDomain(&d); err != nil {
			return nil, err
		}

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		dmn, err := mdl.ToDomain()
		if err != nil {
			return nil, err
		}
		insertedData = append(insertedData, *dmn)
	}

	return insertedData, nil
}

func bootstrapRule(client *postgres.Client) ([]rule.Rule, error) {
	filePath := "./testdata/mock-rule.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []rule.Rule
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []rule.Rule
	for _, d := range data {
		var mdl model.Rule
		if err := mdl.FromDomain(&d); err != nil {
			return nil, err
		}

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		dmn, err := mdl.ToDomain()
		if err != nil {
			return nil, err
		}
		insertedData = append(insertedData, *dmn)
	}

	return insertedData, nil
}

func bootstrapSubscription(client *postgres.Client) ([]subscription.Subscription, error) {
	filePath := "./testdata/mock-subscription.json"
	testFixtureJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []subscription.Subscription
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	var insertedData []subscription.Subscription
	for _, d := range data {
		var mdl model.Subscription
		mdl.FromDomain(&d)

		result := client.GetDB(context.TODO()).Create(&mdl)
		if result.Error != nil {
			return nil, result.Error
		}

		insertedData = append(insertedData, *mdl.ToDomain())
	}

	return insertedData, nil
}

func TestMigration(t *testing.T) {
	t.Run("successfully migrate if there is no problem", func(t *testing.T) {
		logger := log.NewZap()
		client, pool, resource, err := newTestClient(logger)
		if err != nil {
			t.Fatal(err)
		}

		if err = client.Migrate(); err != nil {
			t.Fatal(err)
		}

		if err = purgeDocker(pool, resource); err != nil {
			t.Fatal(err)
		}
	})

}

func TestLogs(t *testing.T) {
	wrongPGConfig := postgres.Config{
		Host:     "localhost",
		User:     "test_user",
		Password: "test_pass",
		Name:     "test_db",
		SSLMode:  "disable",
	}

	t.Run("failed migrate if there is a problem with log level error", func(t *testing.T) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		logger := log.NewZap(log.ZapWithConfig(zapConfig))

		_, err := postgres.NewClient(logger, wrongPGConfig)
		if err == nil {
			t.Fatal("should throw error")
		}
	})

	t.Run("failed migrate if there is a problem with log level warn", func(t *testing.T) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		logger := log.NewZap(log.ZapWithConfig(zapConfig))

		_, err := postgres.NewClient(logger, wrongPGConfig)
		if err == nil {
			t.Fatal("should throw error")
		}
	})

	t.Run("failed migrate if there is a problem with log level info", func(t *testing.T) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger := log.NewZap(log.ZapWithConfig(zapConfig))

		_, err := postgres.NewClient(logger, wrongPGConfig)
		if err == nil {
			t.Fatal("should throw error")
		}
	})

	t.Run("failed migrate if there is a problem with log level debug", func(t *testing.T) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger := log.NewZap(log.ZapWithConfig(zapConfig))

		_, err := postgres.NewClient(logger, wrongPGConfig)
		if err == nil {
			t.Fatal("should throw error")
		}
	})

	t.Run("failed migrate if there is a problem with other log level", func(t *testing.T) {
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
		logger := log.NewZap(log.ZapWithConfig(zapConfig))

		_, err := postgres.NewClient(logger, wrongPGConfig)
		if err == nil {
			t.Fatal("should throw error")
		}
	})
}
