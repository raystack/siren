package rules

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/templates"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

// AnyTime is used to expect arbitrary time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository RuleRepository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = NewRepository(db)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestUpsert() {
	s.Run("should insert rule merged with defaults and call cortex APIs", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedTemplateRowsSecondQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnRows(expectedTemplateRowsSecondQuery)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
	})

	s.Run("should update rule merged with defaults and call cortex APIs", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"entity"=$4,"group_name"=$5,"template"=$6,"status"=$7,"variables"=$8 WHERE id = $9`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedTemplateRowsSecondQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables, expectedRule.ID).WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnRows(expectedTemplateRowsSecondQuery)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
	})

	s.Run("should rollback update if cortex API call fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"entity"=$4,"group_name"=$5,"template"=$6,"status"=$7,"variables"=$8 WHERE id = $9`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedTemplateRowsSecondQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables, expectedRule.ID).WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnRows(expectedTemplateRowsSecondQuery)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback insert if cortex API call fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedTemplateRowsSecondQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnRows(expectedTemplateRowsSecondQuery)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if insert query fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if first select query fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if second select query fails", func() {
		mockClient := &cortexCallerMock{}
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should disable alerts if no error from cortex", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(nil)

		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		mockClient.AssertCalled(s.T(), "DeleteRuleGroup", mock.Anything, "foo", "bar")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if delete rule group call fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(errors.New("random error"))

		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should handle deletion of non-existent rule group", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(errors.New("requested resource not found"))

		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if template get query fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
		}

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnError(errors.New("random error"))
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if template get query fails while rendering", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if rule variables json unmarshalling fails", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "test",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `{}`,
		}

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "json: cannot unmarshal object into Go value of type []domain.RuleVariable")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if rule body yaml unmarshalling fail", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		secondSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","entity","group_name","template","status","variables") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "abcd",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedTemplateRowsSecondQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(secondSelectTemplateQuery).WillReturnRows(expectedTemplateRowsSecondQuery)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.EqualError(err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `abcd` into []rulefmt.RuleNode")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should store disabled alerts", func() {
		mockClient := &cortexCallerMock{}
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(nil)
		firstSelectTemplateQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'tmpl'`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_gojek_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND entity = 'gojek' AND group_name = 'bar'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"entity"=$4,"group_name"=$5,"template"=$6,"status"=$7,"variables"=$8 WHERE id = $9`)
		expectedTemplate := &templates.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
			Tags:      []string{"baz"},
			Variables: `[{"name":"for","type":"string","default":"10m","description":"test"},{"name":"team","type":"string","default":"gojek","description":"test"}]`,
		}
		input := &Rule{
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "disabled",
			Variables: `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRule.ID, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.Entity, expectedRule.GroupName, expectedRule.Template,
				expectedRule.Status, expectedRule.Variables)

		expectedTemplateRowsFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectTemplateQuery).WillReturnRows(expectedTemplateRowsFirstQuery)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.Entity, expectedRule.GroupName, expectedRule.Template, expectedRule.Status,
			expectedRule.Variables, expectedRule.ID).WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func (s *RepositoryTestSuite) TestGet() {
	s.Run("should get rules filtered on parameters", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules WHERE namespace = 'test-namespace' AND entity = 'test-entity' AND group_name = 'test-group' AND status = 'test-enabled' AND template = 'test-template'`)
		expectedRules := []Rule{{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].Entity, expectedRules[0].GroupName, expectedRules[0].Template,
				expectedRules[0].Status, expectedRules[0].Variables)

		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

		actualRules, err := s.repository.Get("test-namespace", "test-entity", "test-group", "test-enabled", "test-template")
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should get rules all rules if empty filters passes", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules`)
		expectedRules := []Rule{{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "siren_api_gojek_foo_bar_tmpl",
			Namespace: "foo",
			GroupName: "bar",
			Entity:    "gojek",
			Template:  "tmpl",
			Status:    "enabled",
			Variables: `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "entity", "group_name", "template", "status", "variables"}).
			AddRow(expectedRules[0].ID, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].Entity, expectedRules[0].GroupName, expectedRules[0].Template,
				expectedRules[0].Status, expectedRules[0].Variables)

		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

		actualRules, err := s.repository.Get("", "", "", "", "")
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if any", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules`)
		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnError(errors.New("random error"))
		actualRules, err := s.repository.Get("", "", "", "", "")
		s.EqualError(err, "random error; random error")
		s.Nil(actualRules)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
