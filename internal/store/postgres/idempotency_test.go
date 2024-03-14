package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"

	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/pgc"
)

type IdempotencyRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.IdempotencyRepository
}

func (s *IdempotencyRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	dpg, err := dockertestx.CreatePostgres(
		dockertestx.PostgresWithDetail(
			pgUser, pgPass, pgDBName,
		),
		dockertestx.PostgresWithVersionTag("13"),
	)
	if err != nil {
		s.T().Fatal(err)
	}

	s.pool = dpg.GetPool()
	s.resource = dpg.GetResource()

	dbConfig.URL = dpg.GetExternalConnString()
	s.client, err = pgc.NewClient(logger, dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx = context.TODO()
	s.Require().NoError(migrate(s.ctx, logger, s.client, dbConfig))
	s.repository = postgres.NewIdempotencyRepository(s.client)

	_, err = bootstrapIdempotency(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *IdempotencyRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *IdempotencyRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *IdempotencyRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE idempotencies RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *IdempotencyRepositoryTestSuite) TestCreate() {
	data1 := &notification.Idempotency{
		Scope:          "a-scope",
		Key:            "key-1",
		NotificationID: "1234-5678-9087",
	}
	s.Run("should return the created idempotency data if not exist", func() {
		res, err := s.repository.Create(context.Background(), data1.Scope, data1.Key, data1.NotificationID)
		s.Assert().NoError(err)
		data1 = res
		s.Assert().Equal("a-scope", data1.Scope)
		s.Assert().Equal("key-1", data1.Key)
		s.Assert().Equal("1234-5678-9087", data1.NotificationID)
	})

	s.Run("should return error if scope or key are empty when creating data", func() {
		_, err := s.repository.Create(context.Background(), "", "", "")
		s.Assert().EqualError(err, "scope or key cannot be empty")
	})
}

func (s *IdempotencyRepositoryTestSuite) TestCheck() {
	s.Run("should return the idempotency if exist", func() {
		res, err := s.repository.Check(context.Background(), "test-default", "xxx-yyy")
		s.Assert().NoError(err)

		s.Assert().Equal("test-default", res.Scope)
		s.Assert().Equal("xxx-yyy", res.Key)
		s.Assert().Equal("1234-5678", res.NotificationID)
	})

	s.Run("should return error not found if idempotency not exist", func() {
		_, err := s.repository.Check(context.Background(), "test-default", "random")
		s.Assert().EqualError(err, "requested entity not found")
	})
}

func (s *IdempotencyRepositoryTestSuite) TestDelete() {
	_, err := s.client.ExecContext(s.ctx, `
		INSERT INTO idempotencies (scope, key, created_at, updated_at) VALUES
			('a-scope', 'old-key-1', now() - interval '2 days', now() - interval '2 days'),
			('a-scope', 'old-key-2', now() - interval '2 days', now() - interval '2 days'),
			('a-scope', 'key-1', now(), now())
	`)
	s.Require().Nil(err)

	s.Run("should return not found if no rows outside TTL", func() {
		err := s.repository.Delete(s.ctx, notification.IdempotencyFilter{
			TTL: time.Hour * time.Duration(10000),
		})
		s.Assert().EqualError(err, errors.ErrNotFound.Error())
	})

	s.Run("should remove all idempotencies that are outside TTL", func() {
		err := s.repository.Delete(s.ctx, notification.IdempotencyFilter{
			TTL: time.Second * time.Duration(60),
		})
		s.Assert().Nil(err)

		numRows := 0
		rows, err := s.client.QueryxContext(s.ctx, `
			SELECT * FROM idempotencies
		`)
		s.Require().Nil(err)

		for rows.Next() {
			numRows += 1
		}

		s.Assert().Equal(1, numRows)
	})
}

func TestIdempotencyRepository(t *testing.T) {
	suite.Run(t, new(IdempotencyRepositoryTestSuite))
}
