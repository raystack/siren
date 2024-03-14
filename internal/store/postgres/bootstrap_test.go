package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goto/salt/db"
	saltlog "github.com/goto/salt/log"
	"github.com/goto/siren/core/alert"
	"github.com/goto/siren/core/log"
	"github.com/goto/siren/core/namespace"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/core/provider"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/core/subscription"
	"github.com/goto/siren/core/template"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
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

func migrate(ctx context.Context, logger saltlog.Logger, client *pgc.Client, dbConf db.Config) error {
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

	insertedData, err := repo.List(context.Background(), receiver.Filter{Expanded: true})
	if err != nil {
		return nil, err
	}

	return insertedData, nil
}

func bootstrapAlert(client *pgc.Client) ([]alert.Alert, error) {
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

	var createdAlerts []alert.Alert
	for _, d := range data {
		alrt, err := repo.Create(context.Background(), d)
		if err != nil {
			return nil, err
		}
		createdAlerts = append(createdAlerts, alrt)
	}

	return createdAlerts, nil
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

func bootstrapNotification(client *pgc.Client) ([]notification.Notification, error) {
	filePath := "./testdata/mock-notification.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []notification.Notification
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewNotificationRepository(client)

	var insertedData []notification.Notification
	for _, d := range data {
		newD, err := repo.Create(context.Background(), d)
		if err != nil {
			return nil, err
		}

		insertedData = append(insertedData, newD)
	}

	return insertedData, nil
}

func bootstrapSilence(client *pgc.Client) ([]string, error) {
	filePath := "./testdata/mock-silence.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []silence.Silence
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewSilenceRepository(client)

	var silenceIDs []string
	for _, d := range data {
		id, err := repo.Create(context.Background(), d)
		if err != nil {
			return nil, err
		}

		silenceIDs = append(silenceIDs, id)
	}

	return silenceIDs, nil
}

func bootstrapNotificationLog(
	client *pgc.Client,
	namespaces []namespace.EncryptedNamespace,
	subscriptions []subscription.Subscription,
	receivers []receiver.Receiver,
	silenceIDs []string,
	notifications []notification.Notification,
	alerts []alert.Alert,
) error {
	var data = []log.Notification{
		{
			NamespaceID:    namespaces[0].ID,
			NotificationID: notifications[0].ID,
			SubscriptionID: subscriptions[0].ID,
			ReceiverID:     receivers[0].ID,
			AlertIDs:       []int64{int64(alerts[0].ID)},
			SilenceIDs:     []string{silenceIDs[0], silenceIDs[1]},
		},
		{
			NamespaceID:    namespaces[1].ID,
			NotificationID: notifications[1].ID,
			SubscriptionID: subscriptions[1].ID,
			ReceiverID:     receivers[1].ID,
			AlertIDs:       []int64{int64(alerts[1].ID)},
			SilenceIDs:     []string{silenceIDs[1]},
		},
		{
			NamespaceID:    namespaces[2].ID,
			NotificationID: notifications[1].ID,
			SubscriptionID: subscriptions[2].ID,
			AlertIDs:       []int64{int64(alerts[0].ID), int64(alerts[2].ID)},
			SilenceIDs:     []string{silenceIDs[0]},
		},
		{
			NamespaceID:    namespaces[2].ID,
			NotificationID: notifications[1].ID,
			SubscriptionID: subscriptions[2].ID,
			AlertIDs:       []int64{int64(alerts[0].ID), int64(alerts[2].ID)},
		},
	}

	repo := postgres.NewLogRepository(client)

	if err := repo.BulkCreate(context.Background(), data); err != nil {
		return err
	}

	return nil
}

func bootstrapIdempotency(client *pgc.Client) ([]notification.Idempotency, error) {
	filePath := "./testdata/mock-idempotency.json"
	testFixtureJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data []notification.Idempotency
	if err = json.Unmarshal(testFixtureJSON, &data); err != nil {
		return nil, err
	}

	repo := postgres.NewIdempotencyRepository(client)

	var insertedData []notification.Idempotency
	for _, d := range data {
		newD, err := repo.Create(context.Background(), d.Scope, d.Key, d.NotificationID)
		if err != nil {
			return nil, err
		}

		insertedData = append(insertedData, *newD)
	}

	return insertedData, nil
}
