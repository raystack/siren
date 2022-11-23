package postgres_test

import (
	"context"
	"fmt"

	"github.com/odpf/salt/db"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/internal/store/postgres"
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
