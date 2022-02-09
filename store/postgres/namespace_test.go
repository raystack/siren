package postgres

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/store/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type NamespaceRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository model.NamespaceRepository
}

func (s *NamespaceRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = NewRepository(db)
}

func (s *NamespaceRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *NamespaceRepositoryTestSuite) TestList() {
	s.Run("should get all namespaces", func() {
		expectedQuery := regexp.QuoteMeta(`select * from namespaces`)
		labels := make(model.StringStringMap)
		labels["foo"] = "bar"

		namespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		expectedNamespaces := []*model.Namespace{namespace}

		expectedRows := sqlmock.NewRows([]string{"id", "provider_id", "urn", "name", "credentials", "labels", "created_at", "updated_at"}).
			AddRow(namespace.Id, namespace.ProviderId, namespace.Urn, namespace.Name, namespace.Credentials,
				json.RawMessage(`{"foo": "bar"}`), namespace.CreatedAt, namespace.UpdatedAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualNamespaces, err := s.repository.List()
		s.Equal(expectedNamespaces, actualNamespaces)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedQuery := regexp.QuoteMeta(`select * from namespaces`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualNamespaces, err := s.repository.List()
		s.Nil(actualNamespaces)
		s.EqualError(err, "random error")
	})
}

func (s *NamespaceRepositoryTestSuite) TestCreate() {
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should create a namespace", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "namespaces" 
											("provider_id","urn","name","credentials","labels","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedNamespace.ProviderId, expectedNamespace.Urn,
			expectedNamespace.Name, expectedNamespace.Credentials, expectedNamespace.Labels,
			expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt, expectedNamespace.Id).
			WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Name, expectedNamespace.ProviderId, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt,
				expectedNamespace.Id)

		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualNamespace, err := s.repository.Create(expectedNamespace)
		s.Equal(expectedNamespace, actualNamespace)
		s.Nil(err)
	})

	s.Run("should return errors in creating a namespace", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "namespaces" 
											("provider_id","urn","name","credentials","labels","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedNamespace.ProviderId, expectedNamespace.Urn,
			expectedNamespace.Name, expectedNamespace.Credentials, expectedNamespace.Labels,
			expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt, expectedNamespace.Id).
			WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Create(expectedNamespace)
		s.EqualError(err, "random error")
		s.Nil(actualNamespace)
	})

	s.Run("should return error if finding newly inserted namespace fails", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "namespaces" 
											("provider_id","urn","name","credentials","labels","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedNamespace.ProviderId, expectedNamespace.Urn,
			expectedNamespace.Name, expectedNamespace.Credentials, expectedNamespace.Labels,
			expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt, expectedNamespace.Id).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Create(expectedNamespace)
		s.EqualError(err, "random error")
		s.Nil(actualNamespace)
	})
}

func (s *NamespaceRepositoryTestSuite) TestGet() {
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should get namespace by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Name, expectedNamespace.ProviderId, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt,
				expectedNamespace.Id)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualNamespace, err := s.repository.Get(1)
		s.Equal(expectedNamespace, actualNamespace)
		s.Nil(err)
	})

	s.Run("should return nil if namespaces of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualNamespace, err := s.repository.Get(1)
		s.Nil(actualNamespace)
		s.Nil(err)
	})

	s.Run("should return error in getting namespace of given id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Get(1)
		s.Nil(actualNamespace)
		s.EqualError(err, "random error")
	})
}

func (s *NamespaceRepositoryTestSuite) TestUpdate() {
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should update a namespace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6 
			WHERE id = $7 AND "id" = $8`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		input := &model.Namespace{
			Id:          1,
			ProviderId:  2,
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Name, expectedNamespace.ProviderId, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt,
				expectedNamespace.Id)
		expectedRows2 := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Name, input.ProviderId,
				json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo": "bar"}`),
				expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt, expectedNamespace.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.ProviderId, input.Name, input.Credentials, input.Labels,
			AnyTime{}, AnyTime{}, input.Id, input.Id).WillReturnResult(sqlmock.NewResult(1, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)

		actualNamespace, err := s.repository.Update(input)
		s.Equal(uint64(2), actualNamespace.ProviderId)
		s.Nil(err)
	})

	s.Run("should return error if namespace does not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		input := &model.Namespace{
			Id:          1,
			ProviderId:  2,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualNamespace, err := s.repository.Update(input)
		s.Nil(actualNamespace)
		s.EqualError(err, "namespace doesn't exist")
	})

	s.Run("should return error in finding the namespace", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		input := &model.Namespace{
			Id:          1,
			ProviderId:  2,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Update(input)
		s.Nil(actualNamespace)
		s.EqualError(err, "random error")
	})

	s.Run("should return error updating the provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6 
			WHERE id = $7 AND "id" = $8`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		input := &model.Namespace{
			Id:          1,
			ProviderId:  2,
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Urn, expectedNamespace.Name, expectedNamespace.ProviderId,
				json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt,
				expectedNamespace.UpdatedAt, expectedNamespace.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.ProviderId, input.Name, input.Credentials, input.Labels,
			AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Update(input)
		s.Nil(actualNamespace)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated provider", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6 
			WHERE id = $7 AND "id" = $8`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &model.Namespace{
			Id:          1,
			ProviderId:  1,
			Urn:         "foo",
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		input := &model.Namespace{
			Id:          1,
			ProviderId:  2,
			Name:        "foo",
			Credentials: `{"foo":"bar"}`,
			Labels:      labels,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.Name, expectedNamespace.ProviderId, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt,
				expectedNamespace.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.ProviderId, input.Name, input.Credentials, input.Labels,
			AnyTime{}, AnyTime{}, input.Id, input.Id).WillReturnResult(sqlmock.NewResult(1, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))

		actualNamespace, err := s.repository.Update(input)
		s.Nil(actualNamespace)
		s.EqualError(err, "random error")
	})
}

func (s *NamespaceRepositoryTestSuite) TestDelete() {
	s.Run("should delete namespace of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "namespaces" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repository.Delete(1)
		s.Nil(err)
	})

	s.Run("should return error in deleting namespace of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "namespaces" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnError(errors.New("random error"))

		err := s.repository.Delete(1)
		s.EqualError(err, "random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(NamespaceRepositoryTestSuite))
}
