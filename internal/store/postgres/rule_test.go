package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RuleRepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	repository *postgres.RuleRepository
}

func (s *RuleRepositoryTestSuite) SetupSuite() {
	var err error

	logger := log.NewZap()
	s.db, s.pool, s.resource, err = newTestClient(logger)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	s.repository = postgres.NewRuleRepository(s.db)

	_, err = bootstrapProvider(s.db)
	if err != nil {
		s.T().Fatal(err)
	}

	_, err = bootstrapNamespace(s.db)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RuleRepositoryTestSuite) SetupTest() {
	var err error
	_, err = bootstrapRule(s.db)
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
	return execQueries(context.TODO(), s.db, queries)
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
							Value: "odpf-app",
						},
						{
							Name:  "name",
							Value: "namespace-urn-1",
						},
						{
							Name:  "entity",
							Value: "odpf",
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
							Value: "odpf-web",
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
				Namespace: "namespace-urn-2",
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
							Value: "odpf-web",
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
		Description   string
		RuleToUpsert  *rule.Rule
		PostProcessFn func() error
		ExpectedID    uint64
		ErrString     string
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
			PostProcessFn: func() error { return nil },
			ExpectedID:    uint64(3), // autoincrement in db side
		},
		{
			Description: "should rollback when create a rule but post process return error",
			RuleToUpsert: &rule.Rule{
				Name:              "new-rule-name",
				Enabled:           true,
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				Variables:         []rule.RuleVariable{},
				ProviderNamespace: 1,
			},
			PostProcessFn: func() error { return errors.New("rollback error") },
			ErrString:     "rollback error",
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
			PostProcessFn: func() error { return nil },
			ExpectedID:    2,
		},
		{
			Description: "should conflict if all unique components are same and name different",
			RuleToUpsert: &rule.Rule{
				Name:              "the-new-conflict",
				GroupName:         "group-name-1",
				Namespace:         "namespace-urn-2",
				Template:          "template-name-1",
				ProviderNamespace: 2,
			},
			ErrString: "",
		},
		{
			Description:   "should return error if rule is nil",
			PostProcessFn: func() error { return nil },
			ErrString:     "rule domain is nil",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Description, func() {
			got, err := s.repository.UpsertWithTx(s.ctx, tc.RuleToUpsert, tc.PostProcessFn)
			if tc.ErrString != "" {
				if err.Error() != tc.ErrString {
					s.T().Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedID) {
				s.T().Fatalf("got result %+v, expected was %+v", got, tc.ExpectedID)
			}
		})
	}
}

func TestRuleRepository(t *testing.T) {
	suite.Run(t, new(RuleRepositoryTestSuite))
}
