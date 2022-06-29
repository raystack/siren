package postgres_test

import (
	"database/sql"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/internal/store/model"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/internal/store/postgres/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/suite"
)

type NamespaceRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository *postgres.NamespaceRepository
}

func (s *NamespaceRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = postgres.NewNamespaceRepository(db)
}

func (s *NamespaceRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *NamespaceRepositoryTestSuite) TestList() {
	s.Run("should get all namespaces", func() {
		expectedQuery := regexp.QuoteMeta(`select * from namespaces`)
		labels := make(model.StringStringMap)
		labels["foo"] = "bar"

		ns := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        1,
				Provider:  1,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}
		expectedNamespaces := []*namespace.EncryptedNamespace{ns}

		expectedRows := sqlmock.NewRows([]string{"id", "provider_id", "urn", "name", "credentials", "labels", "created_at", "updated_at"}).
			AddRow(ns.ID, ns.Provider, ns.URN, ns.Name, ns.Credentials,
				json.RawMessage(`{"foo": "bar"}`), ns.CreatedAt, ns.UpdatedAt)
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

	insertQuery := regexp.QuoteMeta(`INSERT INTO "namespaces" 
										("provider_id","urn","name","credentials","labels","created_at","updated_at")
										VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
	timeNow := time.Now()

	s.Run("should create a namespace", func() {
		expectedID := uint64(1)
		expectedNamespace := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        expectedID,
				Provider:  1,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			Credentials: `{"foo":"bar"}`,
		}

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				Provider:  1,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			Credentials: `{"foo":"bar"}`,
		}
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(
				expectedNamespace.Provider,
				expectedNamespace.URN,
				expectedNamespace.Name,
				expectedNamespace.Credentials,
				labels,
				expectedNamespace.CreatedAt,
				expectedNamespace.UpdatedAt,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
		s.dbmock.ExpectCommit()

		err := s.repository.Create(input)
		s.Equal(expectedNamespace, input)
		s.Nil(err)
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return errors in creating a namespace", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "namespaces" 
											("provider_id","urn","name","credentials","labels","created_at","updated_at")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				Provider:  1,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(
				input.Provider,
				input.URN,
				input.Name,
				input.Credentials,
				labels,
				input.CreatedAt,
				input.UpdatedAt,
			).
			WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		err := s.repository.Create(input)
		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})
}

func (s *NamespaceRepositoryTestSuite) TestGet() {
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should get namespace by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)
		expectedNamespace := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        1,
				Provider:  1,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		expectedRows := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(expectedNamespace.URN, expectedNamespace.Name, expectedNamespace.Provider, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), expectedNamespace.CreatedAt, expectedNamespace.UpdatedAt,
				expectedNamespace.ID)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualNamespace, err := s.repository.Get(1)
		s.Equal(expectedNamespace, actualNamespace)
		s.Nil(err)
	})

	s.Run("should return not found error if namespaces of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualNamespace, err := s.repository.Get(1)
		s.Nil(actualNamespace)
		s.EqualError(err, "namespace with id 1 not found")
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
	labels["foo"] = "bar-label"

	s.Run("should update a namespace", func() {
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6 
			WHERE id = $7 AND "id" = $8`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1 AND "namespaces"."id" = $1`)

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        1,
				Provider:  2,
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Provider, input.Name, input.Credentials, labels,
			AnyTime{}, AnyTime{}, input.ID, input.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		expectedRows := sqlmock.
			NewRows([]string{"urn", "name", "provider_id", "credentials", "labels", "created_at", "updated_at", "id"}).
			AddRow(input.URN, input.Name, input.Provider, json.RawMessage(`{"foo":"bar"}`),
				json.RawMessage(`{"foo": "bar"}`), input.CreatedAt, input.UpdatedAt,
				input.ID)
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.ID).WillReturnRows(expectedRows)
		s.dbmock.ExpectCommit()

		err := s.repository.Update(input)
		s.Nil(err)
		s.Equal(uint64(2), input.Provider)
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error if namespace does not exist", func() {
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"urn"=$2,"name"=$3,"credentials"=$4,"labels"=$5,"created_at"=$6,"updated_at"=$7
			WHERE id = $8 AND "id" = $9`)

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        99,
				Provider:  2,
				URN:       "foo",
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Provider, input.URN, input.Name, input.Credentials, labels,
			AnyTime{}, AnyTime{}, input.ID, input.ID).WillReturnResult(sqlmock.NewResult(0, 0))
		s.dbmock.ExpectRollback()

		err := s.repository.Update(input)
		s.EqualError(err, "namespace with id 99 not found")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error updating the namespace", func() {
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6
			WHERE id = $7 AND "id" = $8`)

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        1,
				Provider:  2,
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Provider, input.Name, input.Credentials, labels,
			AnyTime{}, AnyTime{}, input.ID, input.ID).
			WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		err := s.repository.Update(input)

		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
	})

	s.Run("should return error getting new namespace", func() {
		updateQuery := regexp.QuoteMeta(`UPDATE "namespaces"
			SET "provider_id"=$1,"name"=$2,"credentials"=$3,"labels"=$4,"created_at"=$5,"updated_at"=$6
			WHERE id = $7 AND "id" = $8`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "namespaces" WHERE id = 1 AND "namespaces"."id" = $1`)

		input := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				ID:        1,
				Provider:  2,
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `{"foo":"bar"}`,
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Provider, input.Name, input.Credentials, labels,
			AnyTime{}, AnyTime{}, input.ID, input.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.ID).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		err := s.repository.Update(input)

		s.EqualError(err, "random error")
		s.Nil(s.dbmock.ExpectationsWereMet())
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
