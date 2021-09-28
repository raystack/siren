package workspace

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
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
	repository WorkspaceRepository
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

func (s *RepositoryTestSuite) TestList() {
	s.Run("should get all workspaces", func() {
		expectedQuery := regexp.QuoteMeta(`select * from workspaces`)
		workspace := Workspace{
			Id:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}
		expectedWorkspaces := []Workspace{workspace}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(workspace.Id, workspace.CreatedAt, workspace.UpdatedAt, workspace.Name, workspace.Urn)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualWorkspaces, err := s.repository.List()
		s.Equal(expectedWorkspaces, actualWorkspaces)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedQuery := regexp.QuoteMeta(`select * from workspaces`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualWorkspaces, err := s.repository.List()
		s.Nil(actualWorkspaces)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestCreate() {
	s.Run("should create a workspace", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "workspaces" ("created_at","updated_at","name","urn","id") 
											VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE urn = 'bar'`)
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt,
				expectedWorkspace.UpdatedAt, expectedWorkspace.Name, expectedWorkspace.Urn)
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
			expectedWorkspace.Name, expectedWorkspace.Urn, expectedWorkspace.Id).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualWorkspace, err := s.repository.Create(expectedWorkspace)
		s.Equal(expectedWorkspace, actualWorkspace)
		s.Nil(err)
	})

	s.Run("should return errors in creating a workspace", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "workspaces" ("created_at","updated_at","name","urn","id")
											VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
			expectedWorkspace.Name, expectedWorkspace.Urn, expectedWorkspace.Id).
			WillReturnError(errors.New("random error"))
		actualWorkspace, err := s.repository.Create(expectedWorkspace)
		s.EqualError(err, "random error")
		s.Nil(actualWorkspace)
	})

	s.Run("should return error if finding newly inserted workspace fails", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "workspaces" ("created_at","updated_at","name","urn","id") 
											VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE urn = 'bar'`)
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
			expectedWorkspace.Name, expectedWorkspace.Urn, expectedWorkspace.Id).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		actualWorkspace, err := s.repository.Create(expectedWorkspace)
		s.EqualError(err, "random error")
		s.Nil(actualWorkspace)
	})
}

func (s *RepositoryTestSuite) TestGet() {
	s.Run("should get workspaces by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 1`)
		expectedWorkspace := &Workspace{
			Id:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}
		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
				expectedWorkspace.Name, expectedWorkspace.Urn)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualWorkspace, err := s.repository.Get(1)
		s.Equal(expectedWorkspace, actualWorkspace)
		s.Nil(err)
	})

	s.Run("should return nil if workspaces of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 1`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualWorkspace, err := s.repository.Get(1)
		s.Nil(actualWorkspace)
		s.Nil(err)
	})

	s.Run("should return error in getting workspace of given id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 1`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualWorkspace, err := s.repository.Get(1)
		s.Nil(actualWorkspace)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestUpdate() {
	s.Run("should update a workspace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "workspaces" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"urn"=$4 WHERE id = $5 AND "id" = $6`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		timeNow := time.Now()
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}

		input := &Workspace{
			Id:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Urn:       "baz",
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
				expectedWorkspace.Name, expectedWorkspace.Urn)

		expectedRows2 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
				expectedWorkspace.Name, input.Urn)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedWorkspace.Name, "baz", expectedWorkspace.Id, expectedWorkspace.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)
		actualWorkspace, err := s.repository.Update(input)
		s.Equal("baz", actualWorkspace.Urn)
		s.Nil(err)
	})

	s.Run("should return error if workspace does not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		timeNow := time.Now()
		input := &Workspace{
			Id:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Urn:       "baz",
		}
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		actualWorkspace, err := s.repository.Update(input)
		s.Nil(actualWorkspace)
		s.EqualError(err, "workspace doesn't exist")
	})

	s.Run("should return error in finding the workspace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		timeNow := time.Now()
		input := &Workspace{
			Id:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Urn:       "baz",
		}
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))
		actualWorkspace, err := s.repository.Update(input)
		s.Nil(actualWorkspace)
		s.EqualError(err, "random error")
	})

	s.Run("should return updating in finding the workspace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "workspaces" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"urn"=$4 WHERE id = $5 AND "id" = $6`)
		timeNow := time.Now()
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}

		input := &Workspace{
			Id:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Urn:       "baz",
		}

		expectedRows1 := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
				expectedWorkspace.Name, expectedWorkspace.Urn)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedWorkspace.Name, "baz", expectedWorkspace.Id, expectedWorkspace.Id).
			WillReturnError(errors.New("random error"))
		actualWorkspace, err := s.repository.Update(input)
		s.Nil(actualWorkspace)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated workspace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "workspaces" SET "created_at"=$1,"updated_at"=$2,"name"=$3,"urn"=$4 WHERE id = $5 AND "id" = $6`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "workspaces" WHERE id = 10`)
		timeNow := time.Now()
		expectedWorkspace := &Workspace{
			Id:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "foo",
			Urn:       "bar",
		}

		input := &Workspace{
			Id:        10,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Name:      "foo",
			Urn:       "baz",
		}

		expectedRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name", "urn"}).
			AddRow(expectedWorkspace.Id, expectedWorkspace.CreatedAt, expectedWorkspace.UpdatedAt,
				expectedWorkspace.Name, expectedWorkspace.Urn)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs(AnyTime{},
			AnyTime{}, expectedWorkspace.Name, "baz", expectedWorkspace.Id, expectedWorkspace.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))
		actualWorkspace, err := s.repository.Update(input)
		s.Nil(actualWorkspace)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	s.Run("should delete workspaces of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "workspaces" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repository.Delete(1)
		s.Nil(err)
	})

	s.Run("should return error in deleting workspace of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "workspaces" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnError(errors.New("random error"))

		err := s.repository.Delete(1)
		s.EqualError(err, "random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
