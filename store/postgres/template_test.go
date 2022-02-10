package postgres_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/store"
	"github.com/odpf/siren/store/model"
	"github.com/odpf/siren/store/postgres"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type TemplateRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository store.TemplatesRepository
}

func (s *TemplateRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = postgres.NewTemplateRepository(db)
}

func (s *TemplateRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *TemplateRepositoryTestSuite) TestIndex() {

	s.Run("should get all templates if tag is not passed", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates"`)
		template := model.Template{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		expectedTemplates := []model.Template{template}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(template.ID, template.CreatedAt,
				template.UpdatedAt, template.Name,
				template.Body, template.Tags,
				template.Variables)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualTemplates, err := s.repository.Index("")
		s.Equal(expectedTemplates, actualTemplates)
		s.Nil(err)
	})

	s.Run("should get templates of matching tags", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE tags @>ARRAY[$1]`)
		template := model.Template{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"foo"},
			Variables: `{"name":"test"}`,
		}
		expectedTemplates := []model.Template{template}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(template.ID, template.CreatedAt, template.UpdatedAt, template.Name, template.Body, template.Tags, template.Variables)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualTemplates, err := s.repository.Index("foo")
		s.Equal(expectedTemplates, actualTemplates)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedErrorMessage := "random error"
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates"`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualTemplates, err := s.repository.Index("")
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplates)
	})
}

func (s *TemplateRepositoryTestSuite) TestGetByName() {

	s.Run("should get template by name", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualTemplates, err := s.repository.GetByName("foo")
		s.Equal(expectedTemplate, actualTemplates)
		s.Nil(err)
	})

	s.Run("should return nil if template not found", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualTemplate, err := s.repository.GetByName("foo")
		s.Nil(actualTemplate)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedErrorMessage := "random error"
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualTemplates, err := s.repository.GetByName("foo")
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplates)
	})
}

func (s *TemplateRepositoryTestSuite) TestDelete() {

	s.Run("should delete template by name", func() {
		deleteQuery := regexp.QuoteMeta(`DELETE FROM "templates" WHERE name = $1`)
		s.dbmock.ExpectExec(deleteQuery).WithArgs("foo").WillReturnResult(sqlmock.NewResult(0, 1))
		err := s.repository.Delete("foo")
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedErrorMessage := "random error"
		deleteQuery := regexp.QuoteMeta(`DELETE FROM "templates" WHERE name = $1`)
		s.dbmock.ExpectExec(deleteQuery).WithArgs("foo").WillReturnError(errors.New(expectedErrorMessage))
		err := s.repository.Delete("foo")
		s.Equal(err.Error(), expectedErrorMessage)
	})
}

func (s *TemplateRepositoryTestSuite) TestUpsert() {

	s.Run("should insert template if not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		insertQuery := regexp.QuoteMeta(`INSERT INTO "templates" ("created_at","updated_at","name","body","tags","variables","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedTemplate.CreatedAt,
			expectedTemplate.UpdatedAt, expectedTemplate.Name,
			expectedTemplate.Body, expectedTemplate.Tags,
			expectedTemplate.Variables, expectedTemplate.ID).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows)
		actualTemplate, err := s.repository.Upsert(expectedTemplate)
		s.Equal(expectedTemplate, actualTemplate)
		s.Nil(err)
	})

	s.Run("should update template if exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		updateQuery := regexp.QuoteMeta(`UPDATE "templates" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"body"=$4,"tags"=$5,"variables"=$6 WHERE id = $7`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		timeNow := time.Now()
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		input := &model.Template{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		expectedRows2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedTemplate.Name,
			expectedTemplate.Body, expectedTemplate.Tags,
			expectedTemplate.Variables, expectedTemplate.ID).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)
		actualTemplate, err := s.repository.Upsert(input)
		s.Equal(expectedTemplate, actualTemplate)
		s.Nil(err)
	})

	s.Run("should return error if first select query fails", func() {
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))
		timeNow := time.Now()
		input := &model.Template{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		actualTemplate, err := s.repository.Upsert(input)
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplate)
	})

	s.Run("should return error if insert fails", func() {
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		insertQuery := regexp.QuoteMeta(`INSERT INTO "templates" ("created_at","updated_at","name","body","tags","variables","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		timeNow := time.Now()
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		input := &model.Template{
			ID:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedTemplate.CreatedAt,
			expectedTemplate.UpdatedAt, expectedTemplate.Name,
			expectedTemplate.Body, expectedTemplate.Tags,
			expectedTemplate.Variables, expectedTemplate.ID).WillReturnError(errors.New("random error"))

		actualTemplate, err := s.repository.Upsert(input)
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplate)
	})

	s.Run("should return error if update fails", func() {
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		updateQuery := regexp.QuoteMeta(`UPDATE "templates" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"body"=$4,"tags"=$5,"variables"=$6 WHERE id = $7`)
		timeNow := time.Now()
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		input := &model.Template{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedTemplate.Name,
			expectedTemplate.Body, expectedTemplate.Tags,
			expectedTemplate.Variables, expectedTemplate.ID).
			WillReturnError(errors.New(expectedErrorMessage))

		actualTemplate, err := s.repository.Upsert(input)
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplate)
	})

	s.Run("should return error if second select query fails", func() {
		expectedErrorMessage := "random error"
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		updateQuery := regexp.QuoteMeta(`UPDATE "templates" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"body"=$4,"tags"=$5,"variables"=$6 WHERE id = $7`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		timeNow := time.Now()
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		input := &model.Template{
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedTemplate.Name,
			expectedTemplate.Body, expectedTemplate.Tags,
			expectedTemplate.Variables, expectedTemplate.ID).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New(expectedErrorMessage))
		actualTemplate, err := s.repository.Upsert(input)
		s.Equal(err.Error(), expectedErrorMessage)
		s.Empty(actualTemplate)
	})
}

func (s *TemplateRepositoryTestSuite) TestRender() {

	s.Run("should render template body from the input", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "The quick [[.color]] fox jumped over the [[.adjective]] dog.",
			Tags:      []string{"baz"},
			Variables: `[{"name":"color","default":"brown","type":"string","description":"test"}, {"name":"adjective","default":"lazy","type":"string","description":"test"}]`,
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)
		expectedBody := "The quick red fox jumped over the dumb dog."
		inputBody := make(map[string]string)
		inputBody["color"] = "red"
		inputBody["adjective"] = "dumb"
		renderedBody, err := s.repository.Render("foo", inputBody)
		s.Equal(expectedBody, renderedBody)
		s.Nil(err)
	})

	s.Run("should render template body enriched with defaults", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		expectedTemplate := &model.Template{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "The quick [[.color]] fox jumped over the [[.adjective]] dog.",
			Tags:      []string{"baz"},
			Variables: `[{"name":"color","default":"red","type":"string","description":"test"}, {"name":"adjective","default":"lazy","type":"string","description":"test"}]`,
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
			AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
				expectedTemplate.UpdatedAt, expectedTemplate.Name,
				expectedTemplate.Body, expectedTemplate.Tags,
				expectedTemplate.Variables)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)
		expectedBody := "The quick brown fox jumped over the lazy dog."
		inputBody := make(map[string]string)
		inputBody["color"] = "brown"
		renderedBody, err := s.repository.Render("foo", inputBody)
		s.Equal(expectedBody, renderedBody)
		s.Nil(err)
	})

	s.Run("should return error if template not found", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))
		inputBody := make(map[string]string)
		inputBody["color"] = "brown"
		renderedBody, err := s.repository.Render("foo", inputBody)
		s.Equal(err.Error(), "template not found")
		s.Equal("", renderedBody)
	})
}

func TestTemplateRepository(t *testing.T) {
	suite.Run(t, new(TemplateRepositoryTestSuite))
}
