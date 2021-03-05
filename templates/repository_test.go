package templates_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/templates"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type RepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository *templates.Repository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = templates.NewRepository(db)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *RepositoryTestSuite) TestIndex() {

	s.Run("should get all templates if tag is not passed", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates"`)
		template := templates.Template{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"baz"},
			Variables: `{"name":"foo"}`,
		}
		expectedTemplates := []templates.Template{template}
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
		template := templates.Template{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Body:      "bar",
			Tags:      []string{"foo"},
			Variables: `{"name":"test"}`,
		}
		expectedTemplates := []templates.Template{template}
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

func (s *RepositoryTestSuite) TestGetByName() {

	s.Run("should get template by name", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		expectedTemplate := &templates.Template{
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

func (s *RepositoryTestSuite) TestDelete() {

	s.Run("should delete template by name", func() {
		//selectQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
		//deleteQuery := regexp.QuoteMeta(`DELETE FROM "templates" WHERE "templates"."id" = $1`)
		//expectedTemplate := &templates.Template{
		//	ID:        10,
		//	CreatedAt: time.Now(),
		//	UpdatedAt: time.Now(),
		//	Name:      "foo",
		//	Body:      "bar",
		//	Tags:      []string{"baz"},
		//	Variables: `{"name":"foo"}`,
		//}
		//expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "body", "tags", "variables"}).
		//	AddRow(expectedTemplate.ID, expectedTemplate.CreatedAt,
		//		expectedTemplate.UpdatedAt, expectedTemplate.Name,
		//		expectedTemplate.Body, expectedTemplate.Tags,
		//		expectedTemplate.Variables)
		//s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		//s.dbmock.ExpectQuery(deleteQuery).WillReturnRows(sqlmock.NewRows(nil))
		//err := s.repository.Delete("foo")
		//s.Nil(err)
	})

	//
	//s.Run("should return error if any", func() {
	//	expectedErrorMessage := "random error"
	//	expectedQuery := regexp.QuoteMeta(`SELECT * FROM "templates" WHERE name = 'foo'`)
	//	s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))
	//
	//	actualTemplates, err := s.repository.GetByName("foo")
	//	s.Equal(err.Error(), expectedErrorMessage)
	//	s.Empty(actualTemplates)
	//})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
