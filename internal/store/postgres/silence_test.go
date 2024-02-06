package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/goto/siren/core/silence"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type SilenceRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.SilenceRepository
	silenceIDs []string
}

func (s *SilenceRepositoryTestSuite) SetupSuite() {
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
	s.repository = postgres.NewSilenceRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
	_, err = bootstrapSubscription(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SilenceRepositoryTestSuite) SetupTest() {

	var err error
	s.silenceIDs, err = bootstrapSilence(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *SilenceRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SilenceRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *SilenceRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE silences RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *SilenceRepositoryTestSuite) TestCreate() {
	type testCase struct {
		Description     string
		SilenceToCreate silence.Silence
		ErrString       string
	}

	var testCases = []testCase{
		{
			Description: "should create a silence type subscription",
			SilenceToCreate: silence.Silence{
				NamespaceID: 1,
				Type:        silence.TypeSubscription,
				TargetID:    1,
				TargetExpression: map[string]any{
					"rule": "true",
				},
			},
		},
		{
			Description: "should create a silence type labels",
			SilenceToCreate: silence.Silence{
				NamespaceID: 1,
				Type:        silence.TypeMatchers,
				TargetExpression: map[string]any{
					"key1": "val1",
				},
			},
		},
		{
			Description: "should return error if a silence is invalid",
			SilenceToCreate: silence.Silence{
				NamespaceID: 111,
				Type:        silence.TypeMatchers,
				TargetExpression: map[string]any{
					"key1": "val1",
				},
			},
			ErrString: "foreign key violation [key (namespace_id)=(111) is not present in table \"namespaces\".]",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			_, err := s.repository.Create(s.ctx, tc.SilenceToCreate)
			if err != nil {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *SilenceRepositoryTestSuite) TestList() {
	type testCase struct {
		description     string
		expectedSilence []silence.Silence
		filter          silence.Filter
		errString       string
	}

	var testCases = []testCase{
		{
			description: "should get all silences",
			expectedSilence: []silence.Silence{
				{
					NamespaceID: 1,
					Type:        "subscription",
					TargetID:    2,
					TargetExpression: map[string]any{
						"rule": "",
					},
				},
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
					},
				},
				{
					NamespaceID: 3,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
				},
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
				},
			},
		},
		{
			description: "should return correct silences when filtered with namespace_id",
			filter: silence.Filter{
				NamespaceID: 2,
			},
			expectedSilence: []silence.Silence{
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
					},
				},
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
				},
			},
		},
		{
			description: "should return correct silences when filtered with subscription_id",
			filter: silence.Filter{
				SubscriptionID: 2,
			},
			expectedSilence: []silence.Silence{
				{
					NamespaceID: 1,
					Type:        "subscription",
					TargetID:    2,
					TargetExpression: map[string]any{
						"rule": "",
					},
				},
			},
		},
		{
			description: "should return correct silences when filtered with labels",
			filter: silence.Filter{
				Match: map[string]string{
					"key1": "value1",
				},
			},
			expectedSilence: []silence.Silence{
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
					},
				},
				{
					NamespaceID: 3,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
				},
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
				},
			},
		},
		{
			description: "should return correct silences when filtered with subscription labels",
			filter: silence.Filter{
				SubscriptionMatch: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			expectedSilence: []silence.Silence{
				{
					NamespaceID: 2,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
					},
				},
				{
					NamespaceID: 3,
					Type:        "labels",
					TargetExpression: map[string]any{
						"key1": "value1",
						"key2": "value2",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.description, func() {
			got, err := s.repository.List(s.ctx, tc.filter)
			if tc.errString != "" {
				if err.Error() != tc.errString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.errString)
				}
			}
			if diff := cmp.Diff(got, tc.expectedSilence, cmpopts.IgnoreFields(silence.Silence{},
				"CreatedAt", "ID")); diff != "" {
				s.T().Fatalf("got diff %+v", diff)
			}
		})
	}
}

func (s *SilenceRepositoryTestSuite) TestGet() {
	type testCase struct {
		Description     string
		ID              string
		ExpectedSilence silence.Silence
		ErrString       string
	}

	var testCases = []testCase{
		{

			Description: "should get a silences",
			ID:          s.silenceIDs[0],
			ExpectedSilence: silence.Silence{
				ID:          "",
				NamespaceID: 1,
				Type:        "subscription",
				TargetID:    2,
				TargetExpression: map[string]any{
					"rule": "",
				},
			},
		},
		{

			Description: "should return error not found when id not exist",
			ID:          "not-exist",
			ExpectedSilence: silence.Silence{
				ID:          "",
				NamespaceID: 1,
				Type:        "subscription",
				TargetID:    2,
				TargetExpression: map[string]any{
					"rule": "",
				},
			},
			ErrString: "sql: no rows in result set",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.Get(s.ctx, tc.ID)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			} else if diff := cmp.Diff(got, tc.ExpectedSilence, cmpopts.IgnoreFields(silence.Silence{},
				"CreatedAt", "ID")); diff != "" {
				s.T().Fatalf("got diff %+v", diff)
			}
		})
	}
}

func (s *SilenceRepositoryTestSuite) TestSoftDelete() {
	s.Run("should delete an entry if success", func() {
		id, err := s.repository.Create(s.ctx, silence.Silence{
			NamespaceID: 1,
			Type:        silence.TypeMatchers,
			TargetExpression: map[string]any{
				"key1": "value1",
			},
		})
		s.Require().NoError(err)

		silences, err := s.repository.List(s.ctx, silence.Filter{})
		s.Require().NoError(err)
		s.Require().Equal(5, len(silences))

		err = s.repository.SoftDelete(s.ctx, id)
		s.Assert().NoError(err)

		silences, err = s.repository.List(s.ctx, silence.Filter{})
		s.Require().NoError(err)
		s.Require().Equal(4, len(silences))
	})

	s.Run("should not delete anything and return nil error if id not found", func() {
		silences, err := s.repository.List(s.ctx, silence.Filter{})
		s.Require().NoError(err)
		s.Require().Equal(4, len(silences))

		err = s.repository.SoftDelete(s.ctx, "random")
		s.Assert().NoError(err)

		silences, err = s.repository.List(s.ctx, silence.Filter{})
		s.Require().NoError(err)
		s.Require().Equal(4, len(silences))
	})
}

func TestSilenceRepository(t *testing.T) {
	suite.Run(t, new(SilenceRepositoryTestSuite))
}
