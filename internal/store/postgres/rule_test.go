package postgres_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/core/rule"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/internal/store/postgres/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type RuleRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository *postgres.RuleRepository
}

func (s *RuleRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = postgres.NewRuleRepository(db)
}

func (s *RuleRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func TestRuleRepository(t *testing.T) {
	suite.Run(t, new(RuleRepositoryTestSuite))
}

func (s *RuleRepositoryTestSuite) TestUpsert() {
	timeNow := time.Now()
	updateQuery := regexp.QuoteMeta(`UPDATE "rules" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"namespace"=$4,"group_name"=$5,"template"=$6,"enabled"=$7,"variables"=$8,"provider_namespace"=$9 WHERE name = $10`)
	insertQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
	findQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = $1`)

	theRule := &rule.Rule{
		CreatedAt:         timeNow,
		UpdatedAt:         timeNow,
		Name:              "siren_api_gojek_foo_bar_tmpl",
		Namespace:         "foo",
		GroupName:         "bar",
		Enabled:           true,
		Template:          "tmpl",
		ProviderNamespace: 1,
		Variables: []rule.RuleVariable{
			{Name: "for", Type: "string", Value: "10m", Description: "test"},
			{Name: "team", Type: "string", Value: "gojek", Description: "test"},
		},
	}
	variablesStr := `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`

	s.Run("should update existing rule", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(theRule.ID, theRule.CreatedAt,
				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
				theRule.GroupName, theRule.Template, theRule.Enabled,
				variablesStr, theRule.ProviderNamespace)
		s.dbmock.ExpectQuery(findQuery).
			WillReturnRows(expectedRow)

		ctx := context.Background()
		err := s.repository.Upsert(ctx, theRule)
		s.Nil(err)
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should create new rule", func() {
		theRule := &rule.Rule{
			CreatedAt:         timeNow,
			UpdatedAt:         timeNow,
			Name:              "siren_api_gojek_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           true,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables: []rule.RuleVariable{
				{Name: "for", Type: "string", Value: "10m", Description: "test"},
				{Name: "team", Type: "string", Value: "gojek", Description: "test"},
			},
		}

		s.dbmock.ExpectExec(updateQuery).
			WithArgs(theRule.CreatedAt, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnResult(sqlmock.NewResult(0, 0))
		expectedID := uint64(1)
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedID, theRule.CreatedAt,
				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
				theRule.GroupName, theRule.Template, theRule.Enabled,
				variablesStr, theRule.ProviderNamespace)
		s.dbmock.ExpectQuery(findQuery).
			WillReturnRows(expectedRow)

		ctx := context.Background()
		err := s.repository.Upsert(ctx, theRule)
		s.Nil(err)
		s.Equal(expectedID, theRule.ID)
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should create using transaction", func() {
		theRule := &rule.Rule{
			CreatedAt:         timeNow,
			UpdatedAt:         timeNow,
			Name:              "siren_api_gojek_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           true,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables: []rule.RuleVariable{
				{Name: "for", Type: "string", Value: "10m", Description: "test"},
				{Name: "team", Type: "string", Value: "gojek", Description: "test"},
			},
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(theRule.CreatedAt, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnResult(sqlmock.NewResult(0, 0))
		expectedID := uint64(1)
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
		expectedRow := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedID, theRule.CreatedAt,
				theRule.UpdatedAt, theRule.Name, theRule.Namespace,
				theRule.GroupName, theRule.Template, theRule.Enabled,
				variablesStr, theRule.ProviderNamespace)
		s.dbmock.ExpectQuery(findQuery).
			WillReturnRows(expectedRow)
		s.dbmock.ExpectCommit()

		ctx := context.Background()
		ctx = s.repository.WithTransaction(ctx)
		err := s.repository.Upsert(ctx, theRule)
		commitErr := s.repository.Commit(ctx)
		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(expectedID, theRule.ID)
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error when updating rule", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnError(errors.New("random error"))

		ctx := context.Background()
		err := s.repository.Upsert(ctx, theRule)
		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error when inserting rule", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
			WillReturnError(errors.New("random error"))

		ctx := context.Background()
		err := s.repository.Upsert(ctx, theRule)
		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error when finding new/updated rule", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace, theRule.Name).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(AnyTime{}, AnyTime{}, theRule.Name, theRule.Namespace, theRule.GroupName, theRule.Template, theRule.Enabled, variablesStr, theRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		s.dbmock.ExpectQuery(findQuery).
			WillReturnError(errors.New("random error"))

		ctx := context.Background()
		err := s.repository.Upsert(ctx, theRule)
		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})
}

func (s *RuleRepositoryTestSuite) TestGet() {
	variablesStr := `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`

	expectedRules := []rule.Rule{{
		ID:                10,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Name:              "siren_api_gojek_foo_bar_tmpl",
		Namespace:         "foo",
		GroupName:         "bar",
		Enabled:           true,
		Template:          "tmpl",
		ProviderNamespace: 1,
		Variables: []rule.RuleVariable{
			{Name: "for", Type: "string", Value: "10m", Description: "test"},
			{Name: "team", Type: "string", Value: "gojek", Description: "test"},
		},
	}}

	s.Run("should get rules filtered on parameters", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = $1 AND namespace = $2 AND group_name = $3 AND template = $4 AND provider_namespace = $5`)

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
				variablesStr, expectedRules[0].ProviderNamespace)

		s.dbmock.ExpectQuery(selectRuleQuery).
			WithArgs("test-name", "test-namespace", "test-group", "test-template", 1).
			WillReturnRows(expectedRuleRows)

		ctx := context.Background()
		actualRules, err := s.repository.Get(ctx, "test-name", "test-namespace", "test-group", "test-template", 1)
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should get rules all rules if empty filters passes", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
				variablesStr, expectedRules[0].ProviderNamespace)

		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

		ctx := context.Background()
		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should use transaction", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
				variablesStr, expectedRules[0].ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectCommit()

		ctx := context.Background()
		ctx = s.repository.WithTransaction(ctx)
		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
		s.Nil(s.repository.Commit(ctx))
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if any", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules"`)
		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnError(errors.New("random error"))
		ctx := context.Background()
		actualRules, err := s.repository.Get(ctx, "", "", "", "", 0)
		s.EqualError(err, "random error")
		s.Nil(actualRules)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
