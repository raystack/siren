package rules

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
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
	expectedTemplate := &domain.Template{
		ID:        10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "tmpl",
		Body:      "-\n    alert: Test\n    expr: 'test-expr'\n    for: '[[.for]]'\n    labels: {severity: WARNING, team:  [[.team]] }\n    annotations: {description: 'test'}\n-\n",
		Tags:      []string{"baz"},
		Variables: []domain.Variable{{
			Name:        "for",
			Default:     "10m",
			Description: "test",
			Type:        "string",
		}, {Name: "team",
			Default:     "gojek",
			Description: "test",
			Type:        "string"},
		}}
	var truebool = true
	var falsebool = false

	dummyTemplateBody := "-\n    alert: Test\n    expr: 'test-expr'\n    for: '20m'\n    labels: {severity: WARNING, team: 'gojek' }\n    annotations: {description: 'test'}\n-\n"

	s.Run("should insert rule merged with defaults and call cortex APIs", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			Variables:         `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
			ProviderNamespace: 1,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			Variables:         `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
			ProviderNamespace: 1,
		}

		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}
		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
	})

	s.Run("should update rule merged with defaults and call cortex APIs", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"group_name"=$4,"template"=$5,"enabled"=$6,"variables"=$7,"provider_namespace"=$8 WHERE id = $9 AND "id" = $10`)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}
		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).
			WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
				expectedRule.ProviderNamespace, expectedRule.Id, expectedRule.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
	})

	s.Run("should rollback update if cortex API call fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"group_name"=$4,"template"=$5,"enabled"=$6,"variables"=$7,"provider_namespace"=$8 WHERE id = $9 AND "id" = $10`)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace, expectedRule.Id, expectedRule.Id).WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback insert if cortex API call fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).
			WithArgs(AnyTime{},
				AnyTime{}, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
				expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if namespace select query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if namespace select query returns no result", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "provider not found")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if provider not supported", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "random",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "provider not supported")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if insert query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}

		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}
		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if first select query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if second select query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		mockClient.AssertNotCalled(s.T(), "CreateRuleGroup")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if third select query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
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
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		mockClient.AssertCalled(s.T(), "DeleteRuleGroup", mock.Anything, "foo", "bar")
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if delete rule group call fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(errors.New("random error"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should handle deletion of non-existent rule group", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(errors.New("requested resource not found"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if template get query fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("GetByName", "tmpl").Return(nil, errors.New("random error"))
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}

		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should rollback if template get query fails while rendering", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return("", errors.New("random error"))
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "random error")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if rule variables json unmarshalling fails", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `{}`,
		}

		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "json: cannot unmarshal object into Go value of type []domain.RuleVariable")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if rule body yaml unmarshalling fail", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		badTemplate := &domain.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "tmpl",
			Body:      "abcd",
			Tags:      []string{"baz"},
			Variables: expectedTemplate.Variables,
		}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(badTemplate.Body, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(badTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectRollback()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `abcd` into []rulefmt.RuleNode")
		s.Nil(actualRule)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should store disabled alerts", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(nil)
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		updateRuleQuery := regexp.QuoteMeta(`UPDATE "rules" SET "updated_at"=$1,"name"=$2,"namespace"=$3,"group_name"=$4,"template"=$5,"enabled"=$6,"variables"=$7,"provider_namespace"=$8 WHERE id = $9 AND "id" = $10`)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &falsebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &falsebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"20m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}
		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}

		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)
		expectedRuleRowsInFirstQuery := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(expectedRuleRowsInFirstQuery)
		s.dbmock.ExpectExec(updateRuleQuery).
			WithArgs(AnyTime{}, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
				expectedRule.ProviderNamespace, expectedRule.Id, expectedRule.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should return error if template not found", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(nil, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &truebool,
			ProviderNamespace: 1,
			Variables:         `[{"name":"for", "type":"string", "value":"20m", "description":"test"}]`,
		}
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.EqualError(err, "template not found")
		s.Nil(actualRule)
	})

	s.Run("should insert disabled rule and not call cortex APIs", func() {
		mockClient := &cortexCallerMock{}
		mockTemplateService := &mocks.TemplatesService{}
		mockTemplateService.On("Render", mock.Anything, mock.Anything).Return(dummyTemplateBody, nil)
		mockTemplateService.On("GetByName", "tmpl").Return(expectedTemplate, nil)
		mockClient.On("CreateRuleGroup", mock.Anything, "foo", mock.Anything).Return(nil)
		mockClient.On("DeleteRuleGroup", mock.Anything, "foo", "bar").Return(errors.New("requested resource not found"))
		namespaceQuery := regexp.QuoteMeta(`SELECT namespaces.urn as namespace_urn, providers.urn as provider_urn, providers.type as provider_type FROM "namespaces" RIGHT JOIN providers on providers.id = namespaces.provider_id WHERE namespaces.id = $1`)
		firstSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		secondSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE name = 'siren_api_bar_foo_foo_bar_tmpl'`)
		thirdSelectRuleQuery := regexp.QuoteMeta(`SELECT * FROM "rules" WHERE namespace = 'foo' AND group_name = 'bar' AND provider_namespace = '1'`)
		insertRuleQuery := regexp.QuoteMeta(`INSERT INTO "rules" ("created_at","updated_at","name","namespace","group_name","template","enabled","variables","provider_namespace") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)

		input := &Rule{
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &falsebool,
			Variables:         `[{"name":"for", "type":"string", "value":"10m", "description":"test"}]`,
			ProviderNamespace: 1,
		}
		expectedRule := &Rule{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_bar_foo_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Template:          "tmpl",
			Enabled:           &falsebool,
			Variables:         `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
			ProviderNamespace: 1,
		}

		expectedNamespace := struct {
			Urn  string
			Purn string
			Type string
		}{
			Urn:  "foo",
			Purn: "bar",
			Type: "cortex",
		}
		expectedNamespaceRow := sqlmock.NewRows([]string{"namespace_urn", "provider_urn", "provider_type"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Purn, expectedNamespace.Type)

		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)
		expectedRuleRowsInGroup := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRule.Id, expectedRule.CreatedAt,
				expectedRule.UpdatedAt, expectedRule.Name, expectedRule.Namespace,
				expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled,
				expectedRule.Variables, expectedRule.ProviderNamespace)

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(namespaceQuery).WillReturnRows(expectedNamespaceRow)
		s.dbmock.ExpectQuery(firstSelectRuleQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertRuleQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedRule.Name, expectedRule.Namespace,
			expectedRule.GroupName, expectedRule.Template, expectedRule.Enabled, expectedRule.Variables,
			expectedRule.ProviderNamespace).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectRuleQuery).WillReturnRows(expectedRuleRows)
		s.dbmock.ExpectQuery(thirdSelectRuleQuery).WillReturnRows(expectedRuleRowsInGroup)
		s.dbmock.ExpectCommit()
		actualRule, err := s.repository.Upsert(input, mockClient, mockTemplateService)
		s.Equal(expectedRule, actualRule)
		s.Nil(err)
	})
}

func (s *RepositoryTestSuite) TestGet() {
	var truebool = true
	s.Run("should get rules filtered on parameters", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules WHERE name = 'test-name' AND namespace = 'test-namespace' AND group_name = 'test-group'  AND template = 'test-template'`)
		expectedRules := []Rule{{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_gojek_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRules[0].Id, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
				expectedRules[0].Variables, expectedRules[0].ProviderNamespace)

		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

		actualRules, err := s.repository.Get("test-name", "test-namespace", "test-group", "test-template")
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	s.Run("should get rules all rules if empty filters passes", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules`)
		expectedRules := []Rule{{
			Id:                10,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			Name:              "siren_api_gojek_foo_bar_tmpl",
			Namespace:         "foo",
			GroupName:         "bar",
			Enabled:           &truebool,
			Template:          "tmpl",
			ProviderNamespace: 1,
			Variables:         `[{"name":"for","type":"string","value":"10m","description":"test"},{"name":"team","type":"string","value":"gojek","description":"test"}]`,
		}}
		expectedRuleRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "namespace", "group_name", "template", "enabled", "variables", "provider_namespace"}).
			AddRow(expectedRules[0].Id, expectedRules[0].CreatedAt,
				expectedRules[0].UpdatedAt, expectedRules[0].Name, expectedRules[0].Namespace,
				expectedRules[0].GroupName, expectedRules[0].Template, expectedRules[0].Enabled,
				expectedRules[0].Variables, expectedRules[0].ProviderNamespace)

		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnRows(expectedRuleRows)

		actualRules, err := s.repository.Get("", "", "", "")
		s.Equal(expectedRules, actualRules)
		s.Nil(err)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	s.Run("should return error if any", func() {
		selectRuleQuery := regexp.QuoteMeta(`SELECT * from rules`)
		s.dbmock.ExpectQuery(selectRuleQuery).WillReturnError(errors.New("random error"))
		actualRules, err := s.repository.Get("", "", "", "")
		s.EqualError(err, "random error; random error")
		s.Nil(actualRules)
		if err := s.dbmock.ExpectationsWereMet(); err != nil {
			s.T().Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
