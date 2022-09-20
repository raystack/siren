package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertest"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/postgres"
)

const (
	logLevelDebug = "debug"
	pgUser        = "test_user"
	pgPass        = "test_pass"
	pgDBName      = "test_db"
)

var (
	dbConfig = db.Config{
		Driver:          "pgx",
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifeTime: 1000,
		MaxQueryTimeout: 1000,
	}
)

func purgeDocker(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		return fmt.Errorf("could not purge resource: %w", err)
	}
	return nil
}

func migrate(ctx context.Context, logger log.Logger, client *postgres.Client, dbConf db.Config) (err error) {
	var queries = []string{
		"DROP SCHEMA public CASCADE",
		"CREATE SCHEMA public",
	}

	err = execQueries(ctx, client, queries)
	if err != nil {
		return
	}

	err = postgres.Migrate(dbConf)
	return
}

// ExecQueries is used for executing list of db query
func execQueries(ctx context.Context, client *postgres.Client, queries []string) error {
	for _, query := range queries {
		_, err := client.GetDB(ctx).QueryContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func bootstrapProvider(client *postgres.Client) ([]provider.Provider, error) {
	filePath := "./testdata/mock-provider.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewProviderRepository(client)

	var data []provider.Provider
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	for _, d := range data {
		if err := repo.Create(context.Background(), &d); err != nil {
			return nil, err
		}
	}

	providers, err := repo.List(context.Background(), provider.Filter{})
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func bootstrapNamespace(client *postgres.Client) ([]namespace.EncryptedNamespace, error) {
	filePath := "./testdata/mock-namespace.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewNamespaceRepository(client)

	var data []namespace.Namespace
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	// encryptService, err := secret.New(testEncryptionKey)
	// if err != nil {
	// 	return nil, err
	// }

	for _, d := range data {
		// plainTextCredentials, err := json.Marshal(d.Credentials)
		// if err != nil {
		// 	return nil, err
		// }
		// cipherTextCredentials, err := encryptService.Encrypt(string(plainTextCredentials))
		// if err != nil {
		// 	return nil, err
		// }
		encryptedNS := namespace.EncryptedNamespace{
			Namespace:   &d,
			Credentials: fmt.Sprintf("%+v", d.Credentials),
		}
		if err := repo.Create(context.Background(), &encryptedNS); err != nil {
			return nil, err
		}
	}

	insertedData, err := repo.List(context.Background())
	if err != nil {
		return nil, err
	}

	return insertedData, nil
}

func bootstrapReceiver(client *postgres.Client) ([]receiver.Receiver, error) {
	filePath := "./testdata/mock-receiver.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewReceiverRepository(client)

	var data []receiver.Receiver
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	for _, d := range data {
		if err := repo.Create(context.Background(), &d); err != nil {
			return nil, err
		}
	}

	insertedData, err := repo.List(context.Background(), receiver.Filter{})
	if err != nil {
		return nil, err
	}

	return insertedData, nil
}

func bootstrapAlert(client *postgres.Client) ([]alert.Alert, error) {
	filePath := "./testdata/mock-alert.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []alert.Alert
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewAlertRepository(client)

	var insertedData []alert.Alert
	for _, d := range data {
		alrt, err := repo.Create(context.Background(), &d)
		if err != nil {
			return nil, err
		}

		insertedData = append(insertedData, *alrt)
	}

	return insertedData, nil
}

func bootstrapTemplate(client *postgres.Client) ([]template.Template, error) {
	filePath := "./testdata/mock-template.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []template.Template
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewTemplateRepository(client)

	for _, d := range data {
		if err := repo.Upsert(context.Background(), &d); err != nil {
			return nil, err
		}
	}

	insertedData, err := repo.List(context.Background(), template.Filter{})
	if err != nil {
		return nil, err
	}
	return insertedData, nil
}

func bootstrapRule(client *postgres.Client) ([]rule.Rule, error) {
	filePath := "./testdata/mock-rule.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []rule.Rule
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewRuleRepository(client)

	var insertedData []rule.Rule
	for _, d := range data {
		err := repo.Upsert(context.Background(), &d)
		if err != nil {
			return nil, err
		}

		insertedData = append(insertedData, d)
	}

	return insertedData, nil
}

func bootstrapSubscription(client *postgres.Client) ([]subscription.Subscription, error) {
	filePath := "./testdata/mock-subscription.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []subscription.Subscription
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewSubscriptionRepository(client)

	var insertedData []subscription.Subscription
	for _, d := range data {
		err := repo.Create(context.Background(), &d)
		if err != nil {
			return nil, err
		}

		insertedData = append(insertedData, d)
	}

	return insertedData, nil
}
