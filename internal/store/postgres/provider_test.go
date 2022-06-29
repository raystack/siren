package postgres_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/internal/store/postgres/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/suite"
)

// AnyTime is used to expect arbitrary time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type ProviderRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository *postgres.ProviderRepository
}

func (s *ProviderRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = postgres.NewProviderRepository(db)
}

func (s *ProviderRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *ProviderRepositoryTestSuite) TestList() {
	s.Run("should get all providers", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers"`)
		credentials := make(model.StringInterfaceMap)
		credentials["foo"] = "bar"
		labels := make(model.StringStringMap)
		labels["foo"] = "bar"

		prv := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProviders := []*provider.Provider{prv}

		expectedRows := sqlmock.
			NewRows([]string{"id", "host", "type", "urn", "name", "credentials", "labels", "created_at", "updated_at"}).
			AddRow(prv.ID, prv.Host, prv.Type, prv.URN, prv.Name, json.RawMessage(`{"foo": "bar"}`),
				json.RawMessage(`{"foo": "bar"}`), prv.CreatedAt, prv.UpdatedAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualProviders, err := s.repository.List(map[string]interface{}{})
		s.Equal(expectedProviders, actualProviders)
		s.Nil(err)
	})

	s.Run("should get all providers by filters", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers"`)
		credentials := make(model.StringInterfaceMap)
		credentials["foo"] = "bar"
		labels := make(model.StringStringMap)
		labels["foo"] = "bar"

		prv := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProviders := []*provider.Provider{prv}

		expectedRows := sqlmock.
			NewRows([]string{"id", "host", "type", "urn", "name", "credentials", "labels", "created_at", "updated_at"}).
			AddRow(prv.ID, prv.Host, prv.Type, prv.URN, prv.Name, json.RawMessage(`{"foo": "bar"}`),
				json.RawMessage(`{"foo": "bar"}`), prv.CreatedAt, prv.UpdatedAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualProviders, err := s.repository.List(map[string]interface{}{
			"urn":  "foo",
			"type": "bar",
		})
		s.Equal(expectedProviders, actualProviders)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers"`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualProviders, err := s.repository.List(map[string]interface{}{})
		s.Nil(actualProviders)
		s.EqualError(err, "random error")
	})
}

func (s *ProviderRepositoryTestSuite) TestCreate() {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should create a provider", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","urn","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)

		modelProvider := &model.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProvider := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   modelProvider.CreatedAt,
			UpdatedAt:   modelProvider.UpdatedAt,
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(modelProvider.Host, modelProvider.URN,
			modelProvider.Name, expectedProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			modelProvider.CreatedAt, modelProvider.UpdatedAt, modelProvider.ID).
			WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(modelProvider.Host, modelProvider.URN, modelProvider.Name, modelProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), modelProvider.CreatedAt,
				modelProvider.UpdatedAt, modelProvider.ID)

		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualProvider, err := s.repository.Create(expectedProvider)
		s.Equal(expectedProvider, actualProvider)
		s.Nil(err)
	})

	s.Run("should return errors in creating a provider", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","urn","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		modelProvider := &model.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProvider := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   modelProvider.CreatedAt,
			UpdatedAt:   modelProvider.UpdatedAt,
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(modelProvider.Host, modelProvider.URN,
			modelProvider.Name, modelProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			modelProvider.CreatedAt, modelProvider.UpdatedAt, modelProvider.ID).
			WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Create(expectedProvider)
		s.EqualError(err, "random error")
		s.Nil(actualProvider)
	})

	s.Run("should return error if finding newly inserted provider fails", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "providers" 
											("host","urn","name","type","credentials","labels","created_at","updated_at","id") 
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)

		modelProvider := &model.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedProvider := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   modelProvider.CreatedAt,
			UpdatedAt:   modelProvider.UpdatedAt,
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(modelProvider.Host, modelProvider.URN,
			modelProvider.Name, modelProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			modelProvider.CreatedAt, modelProvider.UpdatedAt, modelProvider.ID).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Create(expectedProvider)
		s.EqualError(err, "random error")
		s.Nil(actualProvider)
	})
}

func (s *ProviderRepositoryTestSuite) TestGet() {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should get provider by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)
		expectedProvider := &provider.Provider{
			ID:          1,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.URN, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.ID)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualProvider, err := s.repository.Get(1)
		s.Equal(expectedProvider, actualProvider)
		s.Nil(err)
	})

	s.Run("should return not found if provider of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualProvider, err := s.repository.Get(1)
		s.Nil(actualProvider)
		s.EqualError(err, "provider with id 1 not found")
	})

	s.Run("should return error in getting provider of given id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Get(1)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})
}

func (s *ProviderRepositoryTestSuite) TestUpdate() {
	credentials := make(model.StringInterfaceMap)
	credentials["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should update a provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "providers" 
						SET "id"=$1,"host"=$2,"name"=$3,"type"=$4,"credentials"=$5,"labels"=$6,"created_at"=$7,"updated_at"=$8 WHERE id = $9 AND "id" = $10`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		modelProvider := &model.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		expectedProvider := &model.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		input := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.URN, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.ID)
		expectedRows2 := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.URN, "baz", expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.ID)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(modelProvider.ID, modelProvider.Host,
			"baz", modelProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			AnyTime{}, AnyTime{}, modelProvider.ID, modelProvider.ID).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)

		actualProvider, err := s.repository.Update(input)
		s.Equal("baz", actualProvider.Name)
		s.Nil(err)
	})

	s.Run("should return error if provider does not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		input := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "provider with id 10 not found")
	})

	s.Run("should return error in finding the provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		input := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
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
						SET "id"=$1,"host"=$2,"name"=$3,"type"=$4,"credentials"=$5,"labels"=$6,"created_at"=$7,"updated_at"=$8 WHERE id = $9 AND "id" = $10`)
		timeNow := time.Now()
		modelProvider := &model.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		expectedProvider := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		input := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.URN, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.ID)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs(modelProvider.ID, modelProvider.Host,
			"baz", modelProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			AnyTime{}, AnyTime{}, modelProvider.ID, modelProvider.ID).
			WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "providers" 
						SET "id"=$1,"host"=$2,"name"=$3,"type"=$4,"credentials"=$5,"labels"=$6,"created_at"=$7,"updated_at"=$8 WHERE id = $9 AND "id" = $10`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "providers" WHERE id = 10`)
		timeNow := time.Now()
		modelProvider := &model.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		expectedProvider := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			URN:         "foo",
			Name:        "foo",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}
		input := &provider.Provider{
			ID:          10,
			Host:        "foo",
			Type:        "bar",
			Name:        "baz",
			Credentials: credentials,
			Labels:      labels,
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		}

		expectedRows := sqlmock.
			NewRows([]string{"host", "urn", "name", "type", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedProvider.Host, expectedProvider.URN, expectedProvider.Name, expectedProvider.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedProvider.CreatedAt,
				expectedProvider.UpdatedAt, expectedProvider.ID)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs(modelProvider.ID, modelProvider.Host,
			"baz", modelProvider.Type, modelProvider.Credentials, modelProvider.Labels,
			AnyTime{}, AnyTime{}, modelProvider.ID, modelProvider.ID).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))

		actualProvider, err := s.repository.Update(input)
		s.Nil(actualProvider)
		s.EqualError(err, "random error")
	})
}

func (s *ProviderRepositoryTestSuite) TestDelete() {
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

func TestProviderRepository(t *testing.T) {
	suite.Run(t, new(ProviderRepositoryTestSuite))
}
