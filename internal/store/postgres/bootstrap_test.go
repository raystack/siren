package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/alert"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/template"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
)

const (
	pgUser   = "test_user"
	pgPass   = "test_pass"
	pgDBName = "test_db"
)

var (
	dbConfig = db.Config{
		Driver:          "postgres",
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

func migrate(ctx context.Context, logger log.Logger, client *pgc.Client, dbConf db.Config) error {
	var queries = []string{
		"DROP SCHEMA public CASCADE",
		"CREATE SCHEMA public",
	}

	if err := execQueries(ctx, client, queries); err != nil {
		return err
	}

	return postgres.Migrate(dbConf)
}

// ExecQueries is used for executing list of db query
func execQueries(ctx context.Context, client *pgc.Client, queries []string) error {
	for _, query := range queries {
		_, err := client.GetDB(ctx).QueryContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func bootstrapProvider(client *pgc.Client) ([]provider.Provider, error) {
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
		if cerr := repo.Create(context.Background(), &d); cerr != nil {
			return nil, cerr
		}
	}

	providers, err := repo.List(context.Background(), provider.Filter{})
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func bootstrapNamespace(client *pgc.Client) ([]namespace.EncryptedNamespace, error) {
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

	for _, d := range data {

		encryptedNS := namespace.EncryptedNamespace{
			Namespace:        &d,
			CredentialString: fmt.Sprintf("%+v", d.Credentials),
		}
		if cerr := repo.Create(context.Background(), &encryptedNS); cerr != nil {
			return nil, cerr
		}
	}

	insertedData, err := repo.List(context.Background())
	if err != nil {
		return nil, err
	}

	return insertedData, nil
}

func bootstrapReceiver(client *pgc.Client) ([]receiver.Receiver, error) {
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
		if cerr := repo.Create(context.Background(), &d); cerr != nil {
			return nil, cerr
		}
	}

	insertedData, err := repo.List(context.Background(), receiver.Filter{})
	if err != nil {
		return nil, err
	}

	return insertedData, nil
}

func bootstrapAlert(client *pgc.Client) error {
	filePath := "./testdata/mock-alert.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var data []alert.Alert
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return err
	}

	repo := postgres.NewAlertRepository(client)

	for _, d := range data {
		_, err := repo.Create(context.Background(), d)
		if err != nil {
			return err
		}
	}

	return nil
}

func bootstrapTemplate(client *pgc.Client) ([]template.Template, error) {
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
		if cerr := repo.Upsert(context.Background(), &d); cerr != nil {
			return nil, cerr
		}
	}

	insertedData, err := repo.List(context.Background(), template.Filter{})
	if err != nil {
		return nil, err
	}
	return insertedData, nil
}

func bootstrapRule(client *pgc.Client) ([]rule.Rule, error) {
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

func bootstrapSubscription(client *pgc.Client) ([]subscription.Subscription, error) {
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
