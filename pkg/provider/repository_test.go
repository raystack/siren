package provider

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
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
	repository ProviderRepository
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
	s.Run("should get all providers", func() {
		expectedQuery := regexp.QuoteMeta(`select * from providers`)

		credentials := make(StringInterfaceMap)
		credentials["foo"] = "bar"

		labels := make(StringStringMap)
		labels["foo"] = "bar"

		provider := &Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProviders := []*Provider{provider}

		expectedRows := sqlmock.
			NewRows([]string{"id", "host", "type", "name", "credentials", "labels", "created_at", "updated_at"}).
			AddRow(provider.Id, provider.Host, provider.Type, provider.Name, json.RawMessage(`{"foo": "bar"}`),
				json.RawMessage(`{"foo": "bar"}`), provider.CreatedAt, provider.UpdatedAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualProviders, err := s.repository.List()
		s.Equal(expectedProviders, actualProviders)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedQuery := regexp.QuoteMeta(`select * from providers`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualProviders, err := s.repository.List()
		s.Nil(actualProviders)
		s.EqualError(err, "random error")
	})
}


func (s *RepositoryTestSuite) TestCreate() {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"

	labels := make(StringStringMap)
	labels["foo"] = "bar"

	s.Run("should create a provider", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		expectedProvider := &Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedProvider.Host,
			expectedProvider.Name, expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			expectedProvider.CreatedAt, expectedProvider.UpdatedAt, expectedProvider.Id).
			WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualProvider, err := s.repository.Create(expectedProvider)
		s.Equal(expectedProvider, actualProvider)
		s.Nil(err)
	})

	s.Run("should return errors in creating a provider", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		expectedProvider := &Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedProvider.Host,
			expectedProvider.Name, expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			expectedProvider.CreatedAt, expectedProvider.UpdatedAt, expectedProvider.Id).
			WillReturnError(errors.New("random error"))
		actualProvider, err := s.repository.Create(expectedProvider)
		s.EqualError(err, "random error")
		s.Nil(actualProvider)
	})

	s.Run("should return error if finding newly inserted provider fails", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		expectedProvider := &Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedProvider.Host,
			expectedProvider.Name, expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			expectedProvider.CreatedAt, expectedProvider.UpdatedAt, expectedProvider.Id).
			WillReturnRows(sqlmock.NewRows(nil))

		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		actualProvider, err := s.repository.Create(expectedProvider)
		s.EqualError(err, "random error")
		s.Nil(actualProvider)
	})
}


func (s *RepositoryTestSuite) TestGet() {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"

	labels := make(StringStringMap)
	labels["foo"] = "bar"

	s.Run("should get provider by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		expectedProvider := &Provider{
			Id:          1,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)
		actualProvider, err := s.repository.Get(1)
		s.Equal(expectedProvider, actualProvider)
		s.Nil(err)
	})

	s.Run("should return nil if providers of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualProvider, err := s.repository.Get(1)
		s.Nil(actualProvider)
		s.Nil(err)
	})

	s.Run("should return error in getting provider of given id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Get(1)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestUpdate() {
	credentials := make(StringInterfaceMap)
	credentials["foo"] = "bar"

	labels := make(StringStringMap)
	labels["foo"] = "bar"

	s.Run("should update a provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "providers" 
						SET "host"=$1,"name"=$2,"type"=$3,"credentials"=$4,"labels"=$5,"created_at"=$6,"updated_at"=$7 
						WHERE id = $8 AND "id" = $9`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		expectedProvider := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		input := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		expectedRows2 := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, "baz", expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(expectedProvider.Host,
			"baz", expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			AnyTime{}, AnyTime{}, expectedProvider.Id, expectedProvider.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)
		actualProvider, err := s.repository.Update(input)
		s.Equal("baz", actualProvider.Name)
		s.Nil(err)
	})

	s.Run("should return error if provider does not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		input := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))
		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "provider doesn't exist")
	})

	s.Run("should return error in finding the provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		input := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))
		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})

	s.Run("should return error updating the provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "providers" 
						SET "host"=$1,"name"=$2,"type"=$3,"credentials"=$4,"labels"=$5,"created_at"=$6,"updated_at"=$7 
						WHERE id = $8 AND "id" = $9`)
		timeNow := time.Now()

		expectedProvider := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		input := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(expectedProvider.Host,
			"baz", expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			AnyTime{}, AnyTime{}, expectedProvider.Id, expectedProvider.Id).
			WillReturnError(errors.New("random error"))
		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "providers" 
						SET "host"=$1,"name"=$2,"type"=$3,"credentials"=$4,"labels"=$5,"created_at"=$6,"updated_at"=$7 
						WHERE id = $8 AND "id" = $9`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		expectedProvider := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		input := &Provider{
			Id:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"host", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.Id)

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(expectedProvider.Host,
			"baz", expectedProvider.Type, expectedProvider.Credentials, expectedProvider.Labels,
			AnyTime{}, AnyTime{}, expectedProvider.Id, expectedProvider.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))
		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	s.Run("should delete provider of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "providers" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repository.Delete(1)
		s.Nil(err)
	})

	s.Run("should return error in deleting provider of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "providers" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnError(errors.New("random error"))

		err := s.repository.Delete(1)
		s.EqualError(err, "random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
