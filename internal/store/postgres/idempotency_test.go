package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertestx"
	"github.com/odpf/salt/log"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"

	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/pkg/errors"
	"github.com/odpf/siren/pkg/pgc"
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
	dbc, err := db.New(dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}

	s.client, err = pgc.NewClient(logger, dbc)
	if err != nil {
		s.T().Fatal(err)
	}
	s.ctx = context.TODO()
	s.Require().NoError(migrate(s.ctx, logger, s.client, dbConfig))
	s.repository = postgres.NewIdempotencyRepository(s.client)
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

func (s *IdempotencyRepositoryTestSuite) TestInsertReturnOnConflict() {
	data1 := &notification.Idempotency{
		Scope: "a-scope",
		Key:   "key-1",
	}
	s.Run("should return the inserted idempotency data if not exist", func() {
		res, err := s.repository.InsertOnConflictReturning(context.Background(), data1.Scope, data1.Key)
		s.Assert().NoError(err)
		data1 = res
		s.Assert().Equal("a-scope", data1.Scope)
		s.Assert().Equal("key-1", data1.Key)
	})

	s.Run("should return the existing conflicted data if exist", func() {
		res, err := s.repository.InsertOnConflictReturning(context.Background(), data1.Scope, data1.Key)
		s.Assert().NoError(err)

		if diff := cmp.Diff(data1, res, cmpopts.IgnoreFields(notification.Idempotency{}, "UpdatedAt")); diff != "" {
			s.T().Error(diff)
		}
	})
}

func (s *IdempotencyRepositoryTestSuite) TestUpdateSuccess() {
	data := &notification.Idempotency{
		Scope: "a-scope",
		Key:   "existing-key-1",
	}
	res, err := s.repository.InsertOnConflictReturning(context.Background(), data.Scope, data.Key)
	s.Require().NoError(err)
	s.Require().Equal(data.Scope, res.Scope)
	s.Require().Equal(data.Key, res.Key)

	type testCase struct {
		Description string
		ID          uint64
		Status      string
		ErrString   string
	}

	var testCases = []testCase{
		{
			Description: "should update existing idempotency if exist",
			ID:          res.ID,
		},
		{
			Description: "should return not found if idempotency not exist",
			ID:          999,
			ErrString:   "requested entity not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.UpdateSuccess(s.ctx, tc.ID, true)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *IdempotencyRepositoryTestSuite) TestDelete() {
	_, err := s.client.ExecContext(s.ctx, "", "", `
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
		rows, err := s.client.QueryxContext(s.ctx, "", "", `
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
