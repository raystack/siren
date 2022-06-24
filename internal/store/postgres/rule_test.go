package postgres_test

import (
	"context"
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
		ExpectedRules []*rule.Rule
		ErrString     string
	}

	var testCases = []testCase{
		{
			Description: "should get all rules",
			ExpectedRules: []*rule.Rule{
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
			ExpectedRules: []*rule.Rule{
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

// func (s *RuleRepositoryTestSuite) TestUpsert() {
// 	timeNow := time.Now()
// 	updateQuery := regexp.QuoteMeta(`UPDATE "rules" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"namespace"=$4,"group_name"=$5,"template"=$6,"enabled"=$7,"variables"=$8,"provider_namespace"=$9 WHERE name = $10`)
// 	insertQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
// 	findQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = $1`)

// 	theRule := &rule.Rule{
// 		CreatedAt:         timeNow,
// 		UpdatedAt:         timeNow,
// 		Name:              "siren_api_gojek_foo_bar_tmpl",
// 		Namespace:         "foo",
// 		GroupName:         "bar",
// 		Enabled:           true,
// 		Template:          "tmpl",
// 		ProviderNamespace: 1,
// 		Variables: []rule.RuleVariable{
// 			{Name: "for", Type: "string", Value: "10m", Description: "test"},
// 			{Name: "team", Type: "string", Value: "gojek", Description: "test"},
// 		},
// 	}
// 	variablesStr := `[{Name:"for","type":"string",Value:"10m","description":"test"},{Name:"team","type":"string",Value:"gojek","description":"test"}]`

// 	s.Run("should update existing rule", func() {
// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnResult(sqlmock.NewResult(1, 1))
// 		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(theRule.ID, theRule.CreatedAt,
// 				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
// 				theRule.GroupName, theRule.Template, theRule.Enabled,
// 				variablesStr, theRule.ProviderNamespace)
// 		s.dbmock.ExpectQuery(findQuery).
// 			WillReturnRows(expectedRow)

// 		ctx := context.Background()
// 		err := s.repository.Upsert(ctx, theRule)
// 		s.Nil(err)
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})

// 	s.Run("should create new rule", func() {
// 		theRule := &rule.Rule{
// 			CreatedAt:         timeNow,
// 			UpdatedAt:         timeNow,
// 			Name:              "siren_api_gojek_foo_bar_tmpl",
// 			Namespace:         "foo",
// 			GroupName:         "bar",
// 			Enabled:           true,
// 			Template:          "tmpl",
// 			ProviderNamespace: 1,
// 			Variables: []rule.RuleVariable{
// 				{Name: "for", Type: "string", Value: "10m", Description: "test"},
// 				{Name: "team", Type: "string", Value: "gojek", Description: "test"},
// 			},
// 		}

// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(theRule.CreatedAt, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnResult(sqlmock.NewResult(0, 0))
// 		expectedID := uint64(1)
// 		s.dbmock.ExpectQuery(insertQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
// 			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
// 		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(expectedID, theRule.CreatedAt,
// 				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
// 				theRule.GroupName, theRule.Template, theRule.Enabled,
// 				variablesStr, theRule.ProviderNamespace)
// 		s.dbmock.ExpectQuery(findQuery).
// 			WillReturnRows(expectedRow)

// 		ctx := context.Background()
// 		err := s.repository.Upsert(ctx, theRule)
// 		s.Nil(err)
// 		s.Equal(expectedID, theRule.ID)
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})

// 	s.Run("should create using transaction", func() {
// 		theRule := &rule.Rule{
// 			CreatedAt:         timeNow,
// 			UpdatedAt:         timeNow,
// 			Name:              "siren_api_gojek_foo_bar_tmpl",
// 			Namespace:         "foo",
// 			GroupName:         "bar",
// 			Enabled:           true,
// 			Template:          "tmpl",
// 			ProviderNamespace: 1,
// 			Variables: []rule.RuleVariable{
// 				{Name: "for", Type: "string", Value: "10m", Description: "test"},
// 				{Name: "team", Type: "string", Value: "gojek", Description: "test"},
// 			},
// 		}

// 		s.dbmock.ExpectBegin()
// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(theRule.CreatedAt, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnResult(sqlmock.NewResult(0, 0))
// 		expectedID := uint64(1)
// 		s.dbmock.ExpectQuery(insertQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
// 			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
// 		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(expectedID, theRule.CreatedAt,
// 				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
// 				theRule.GroupName, theRule.Template, theRule.Enabled,
// 				variablesStr, theRule.ProviderNamespace)
// 		s.dbmock.ExpectQuery(findQuery).
// 			WillReturnRows(expectedRow)
// 		s.dbmock.ExpectCommit()

// 		ctx := context.Background()
// 		ctx = s.repository.WithTransaction(ctx)
// 		err := s.repository.Upsert(ctx, theRule)
// 		commitErr := s.repository.Commit(ctx)
// 		s.Nil(commitErr)
// 		s.Nil(err)
// 		s.Equal(expectedID, theRule.ID)
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})

// 	s.Run("should return error when updating rule", func() {
// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnError(errors.New("random error"))

// 		ctx := context.Background()
// 		err := s.repository.Upsert(ctx, theRule)
// 		s.EqualError(err, "random error")
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})

// 	s.Run("should return error when inserting rule", func() {
// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnResult(sqlmock.NewResult(0, 0))
// 		s.dbmock.ExpectQuery(insertQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
// 			WillReturnError(errors.New("random error"))

// 		ctx := context.Background()
// 		err := s.repository.Upsert(ctx, theRule)
// 		s.EqualError(err, "random error")
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})

// 	s.Run("should return error when finding new/updated rule", func() {
// 		s.dbmock.ExpectExec(updateQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
// 			WillReturnResult(sqlmock.NewResult(0, 0))
// 		s.dbmock.ExpectQuery(insertQuery).
// 			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
// 			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
// 		s.dbmock.ExpectQuery(findQuery).
// 			WillReturnError(errors.New("random error"))

// 		ctx := context.Background()
// 		err := s.repository.Upsert(ctx, theRule)
// 		s.EqualError(err, "random error")
// 		s.Nil(s.dbmock.ExpectationsWereMet())
// 	})
// }

// func (s *RuleRepositoryTestSuite) TestGet() {
// 	variablesStr := `[{Name:"for","type":"string",Value:"10m","description":"test"},{Name:"team","type":"string",Value:"gojek","description":"test"}]`

// 	expectedRules := []rule.Rule{{
// 		ID:                10,
// 		CreatedAt:         time.Now(),
// 		UpdatedAt:         time.Now(),
// 		Name:              "siren_api_gojek_foo_bar_tmpl",
// 		Namespace:         "foo",
// 		GroupName:         "bar",
// 		Enabled:           true,
// 		Template:          "tmpl",
// 		ProviderNamespace: 1,
// 		Variables: []rule.RuleVariable{
// 			{Name: "for", Type: "string", Value: "10m", Description: "test"},
// 			{Name: "team", Type: "string", Value: "gojek", Description: "test"},
// 		},
// 	}}

// 	s.Run("should get rules filtered on parameters", func() {
// 		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = $1 AND namespace = $2 AND group_name = $3 AND template = $4 AND provider_namespace = $5`)

// 		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
// 				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
// 				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
// 				variablesStr, expectedRules[0].ProviderNamespace)

// 		s.dbmock.ExpectQuery(selectRuleQuery).
// 			WithArgs("test-name", "test-namespace", "test-group", "test-template", 1).
// 			WillReturnRows(expectedRuleRows)

// 		ctx := context.Background()
// 		actualRules, err := s.repository.Get(ctx, "test-name", "test-namespace", "test-group", "test-template", 1)
// 		s.Equal(expectedRules, actualRules)
// 		s.Nil(err)
// 		if err := s.dbmock.ExpectationsWereMet(); err != nil {
// 			s.T().Errorf("there were unfulfilled expectations: %s", err)
// 		}
// 	})

// 	s.Run("should get rules all rules if empty filters passes", func() {
// 		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
// 		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
// 				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
// 				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
// 				variablesStr, expectedRules[0].ProviderNamespace)

// 		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

// 		ctx := context.Background()
// 		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
// 		s.Equal(expectedRules, actualRules)
// 		s.Nil(err)
// 		if err := s.dbmock.ExpectationsWereMet(); err != nil {
// 			s.T().Errorf("there were unfulfilled expectations: %s", err)
// 		}
// 	})

// 	s.Run("should use transaction", func() {
// 		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
// 		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
// 			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
// 				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
// 				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
// 				variablesStr, expectedRules[0].ProviderNamespace)

// 		s.dbmock.ExpectBegin()
// 		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)
// 		s.dbmock.ExpectCommit()

// 		ctx := context.Background()
// 		ctx = s.repository.WithTransaction(ctx)
// 		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
// 		s.Nil(s.repository.Commit(ctx))
// 		s.Equal(expectedRules, actualRules)
// 		s.Nil(err)
// 		if err := s.dbmock.ExpectationsWereMet(); err != nil {
// 			s.T().Errorf("there were unfulfilled expectations: %s", err)
// 		}
// 	})

// 	s.Run("should return error if any", func() {
// 		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
// 		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnError(errors.New("random error"))
// 		ctx := context.Background()
// 		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
// 		s.EqualError(err, "random error")
// 		s.Nil(actualRules)
// 		if err := s.dbmock.ExpectationsWereMet(); err != nil {
// 			s.T().Errorf("there were unfulfilled expectations: %s", err)
// 		}
// 	})
// }

func TestRuleRepository(t *testing.T) {
	suite.Run(t, new(RuleRepositoryTestSuite))
}
