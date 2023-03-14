package postgres_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/goto/salt/db"
	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/internal/store/postgres"
	"github.com/goto/siren/pkg/pgc"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type RuleRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	client     *pgc.Client
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.RuleRepository
}

func (s *RuleRepositoryTestSuite) SetupSuite() {
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
	s.repository = postgres.NewRuleRepository(s.client)

	_, err = bootstrapProvider(s.client)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapNamespace(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapRule(s.client)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) TearDownSuite() {
	// Clean tests
	if err := purgeDocker(s.pool, s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) TearDownTest() {
	if err := s.cleanup(); err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) cleanup() error {
	queries := []string{
		"TRUNCATE TABLE rules RESTART IDENTITY CASCADE",
	}
	return execQueries(context.TODO(), s.client, queries)
}

func (s *RuleRepositoryTestSuite) TestList() {
	type testCase struct {
		Description   string
		Filter        rule.Filter
		ExpectedRules []rule.Rule
		ErrString     string
	}

	var testCases = []testCase{
		{
			Description: "should get all rules",
			ExpectedRules: []rule.Rule{
				{
					ID:        1,
					Name:      "prefix_provider-urn-1_namespace-urn-1_namespace-1_group-name-1_template-name-1",
					GroupName: "group-name-1",
					Namespace: "namespace-urn-1",
					Template:  "template-name-1",
					Variables: []rule.RuleVariable{
						{
							Name:  "WARN_THRESHOLD",
							Value: "200",
							Type:  "int",
						},
						{
							Name:  "CRIT_THRESHOLD",
							Value: "300",
							Type:  "int",
						},
						{
							Name:  "team",
							Value: "gotocompany-app",
						},
						{
							Name:  "name",
							Value: "namespace-urn-1",
						},
						{
							Name:  "entity",
							Value: "gotocompany",
						},
					},
					ProviderNamespace: 1,
				},
				{
					ID:        2,
					Name:      "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
					Enabled:   true,
					GroupName: "group-name-1",
					Namespace: "namespace-urn-2",
					Template:  "template-name-1",
					Variables: []rule.RuleVariable{
						{
							Name:  "WARN_THRESHOLD",
							Value: "5000",
						},
						{
							Name:  "CRIT_THRESHOLD",
							Value: "15000",
						},
						{
							Name:  "team",
							Value: "gotocompany-web",
						},
						{
							Name:  "name",
							Value: "namespace-urn-2",
						},
					},
					ProviderNamespace: 2,
				},
			},
		},
		{
			Description: "should get filtered rules",
			Filter: rule.Filter{
				Namespace:    "namespace-urn-2",
				Name:         "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
				TemplateName: "template-name-1",
			},
			ExpectedRules: []rule.Rule{
				{
					ID:        2,
					Name:      "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
					Enabled:   true,
					GroupName: "group-name-1",
					Namespace: "namespace-urn-2",
					Template:  "template-name-1",
					Variables: []rule.RuleVariable{
						{
							Name:  "WARN_THRESHOLD",
							Value: "5000",
						},
						{
							Name:  "CRIT_THRESHOLD",
							Value: "15000",
						},
						{
							Name:  "team",
							Value: "gotocompany-web",
						},
						{
							Name:  "name",
							Value: "namespace-urn-2",
						},
					},
					ProviderNamespace: 2,
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.List(s.ctx, tc.Filter)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedRules, cmpopts.IgnoreFields(rule.Rule{}, "CreatedAt", "UpdatedAt")) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedRules)
			}
		})
	}
}

func (s *RuleRepositoryTestSuite) TestUpsert() {
	type testCase struct {
		Description  string
		RuleToUpsert *rule.Rule
		ExpectedID   uint64
		ErrString    string
	}

	var testCases = []testCase{
		{
			Description: "should create a rule if doesn't exist",
			RuleToUpsert: &rule.Rule{
				Name:              "new-rule-name",
				Enabled:           true,
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 1,
			},
			ExpectedID: uint64(3), // autoincrement in db side
		},
		{
			Description: "should update a rule if already exist",
			RuleToUpsert: &rule.Rule{
				ID:                2,
				Name:              "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
				Enabled:           true,
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 2,
			},
			ExpectedID: 2,
		},
		{
			Description: "should conflict if all unique components are same and name different",
			RuleToUpsert: &rule.Rule{
				ID:                500,
				Name:              "prefix_provider-urn-1_namespace-urn-1_namespace-1_group-name-1_template-name-1",
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				ProviderNamespace: 2,
			},
			ErrString: "rule conflicted with existing",
		},
		{
			Description: "should return foreign key violation if provider namespace does not exist when update",
			RuleToUpsert: &rule.Rule{
				ID:                2,
				Name:              "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
				Enabled:           true,
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 2000,
			},
			ErrString: "provider's namespace does not exist",
		},
		{
			Description: "should return foreign key violation if provider namespace does not exist when insert",
			RuleToUpsert: &rule.Rule{
				Name:              "new-rule-name-2",
				Enabled:           true,
				GroupName:         "foo",
				Namespace:         "bar",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 1000,
			},
			ErrString: "provider's namespace does not exist",
		},
		{
			Description: "should return foreign key violation if provider namespace does not exist when update",
			RuleToUpsert: &rule.Rule{
				ID:                2,
				Name:              "prefix_provider-urn-1_namespace-urn-2_namespace-2_group-name-1_template-name-1",
				Enabled:           true,
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 2000,
			},
			ErrString: "provider's namespace does not exist",
		},
		{
			Description: "should return foreign key violation if provider namespace does not exist when insert",
			RuleToUpsert: &rule.Rule{
				Name:              "new-rule-name-2",
				Enabled:           true,
				GroupName:         "foo",
				Namespace:         "bar",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 1000,
			},
			ErrString: "provider's namespace does not exist",
		},
		{
			Description: "should return error if rule is nil",
			ErrString:   "rule domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			err := s.repository.Upsert(s.ctx, tc.RuleToUpsert)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
		})
	}
}

func (s *RuleRepositoryTestSuite) TestTransaction() {
	s.Run("successfully commit transaction", func() {
		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Upsert(ctx, &rule.Rule{
			ID:                888,
			Name:              "test-commit",
			Enabled:           true,
			GroupName:         "group-name-1",
			Namespace:         "ns-test-commit",
			Template:          "template-name-1",
			Variables:         []rule.RuleVariable{},
			ProviderNamespace: 2,
		})
		s.NoError(err)

		err = s.repository.Commit(ctx)
		s.NoError(err)

		fetchedRules, err := s.repository.List(s.ctx, rule.Filter{
			Name: "test-commit",
		})
		s.NoError(err)
		s.Len(fetchedRules, 1)
	})

	s.Run("successfully rollback transaction", func() {
		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Upsert(ctx, &rule.Rule{
			ID:                999,
			Name:              "test-rollback",
			Enabled:           true,
			GroupName:         "group-name-1",
			Namespace:         "ns-test-rollback",
			Template:          "template-name-1",
			Variables:         []rule.RuleVariable{},
			ProviderNamespace: 2,
		})
		s.NoError(err)

		err = s.repository.Rollback(ctx, nil)
		s.NoError(err)

		fetchedRules, err := s.repository.List(s.ctx, rule.Filter{
			Name: "test-rollback",
		})
		s.NoError(err)
		s.Len(fetchedRules, 0)
	})
}

func TestRuleRepository(t *testing.T) {
	suite.Run(t, new(RuleRepositoryTestSuite))
}
