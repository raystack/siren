package postgres_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/store"
	"github.com/odpf/siren/store/postgres"
	"github.com/stretchr/testify/suite"
)

type ReceiverRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository store.ReceiverRepository
}

func (s *ReceiverRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	repo := postgres.NewReceiverRepository(db)
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = repo
}

func (s *ReceiverRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *ReceiverRepositoryTestSuite) TestList() {
	s.Run("should get all receivers", func() {
		expectedQuery := regexp.QuoteMeta(`select * from receivers`)
		configurations := make(map[string]interface{})
		configurations["foo"] = "bar"
		labels := make(map[string]string)
		labels["foo"] = "bar"

		receiver := &domain.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		expectedReceivers := []*domain.Receiver{receiver}

		expectedRows := sqlmock.
			NewRows([]string{"id", "name", "type", "labels", "configurations", "created_at", "updated_at"}).
			AddRow(receiver.Id, receiver.Name, receiver.Type, json.RawMessage(`{"foo": "bar"}`),
				json.RawMessage(`{"foo": "bar"}`), receiver.CreatedAt, receiver.UpdatedAt)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualReceivers, err := s.repository.List()
		s.Equal(expectedReceivers, actualReceivers)
		s.Nil(err)
	})

	s.Run("should return error if any", func() {
		expectedQuery := regexp.QuoteMeta(`select * from receivers`)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.List()
		s.Nil(actualReceiver)
		s.EqualError(err, "random error")
	})
}

func (s *ReceiverRepositoryTestSuite) TestCreate() {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	s.Run("should create a receiver", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "receivers"
											("name","type","labels","configurations","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		expectedReceiver := &domain.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedReceiver.Name,
			expectedReceiver.Type, json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			expectedReceiver.CreatedAt, expectedReceiver.UpdatedAt, expectedReceiver.Id).
			WillReturnRows(sqlmock.NewRows(nil))

		err := s.repository.Create(expectedReceiver)
		s.Nil(err)
	})

	s.Run("should return errors in creating a receiver", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "receivers"
											("name","type","labels","configurations","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		expectedReceiver := &domain.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedReceiver.Name, expectedReceiver.Type,
			json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			expectedReceiver.CreatedAt, expectedReceiver.UpdatedAt, expectedReceiver.Id).
			WillReturnError(errors.New("random error"))

		err := s.repository.Create(expectedReceiver)
		s.EqualError(err, "random error")
	})
}

func (s *ReceiverRepositoryTestSuite) TestGet() {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	s.Run("should get receiver by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)
		expectedReceiver := &domain.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		expectedRows := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow(expectedReceiver.Name, expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(expectedRows)

		actualReceiver, err := s.repository.Get(1)
		s.Equal(expectedReceiver, actualReceiver)
		s.Nil(err)
	})

	s.Run("should return error if receiver of given id does not exist", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualReceiver, err := s.repository.Get(1)
		s.Nil(actualReceiver)
		s.EqualError(err, "receiver not found: 1")
	})

	s.Run("should return error in getting receiver of given id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)

		s.dbmock.ExpectQuery(expectedQuery).WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Get(1)
		s.Nil(actualReceiver)
		s.EqualError(err, "random error")
	})
}

func (s *ReceiverRepositoryTestSuite) TestUpdate() {
	configurations := make(map[string]interface{})
	configurations["foo"] = "bar"
	labels := make(map[string]string)
	labels["foo"] = "bar"

	updateQuery := regexp.QuoteMeta(`UPDATE "receivers"
						SET "id"=$1,"name"=$2,"type"=$3,"labels"=$4,"configurations"=$5,"created_at"=$6,"updated_at"=$7
						WHERE id = $8 AND "id" = $9`)
	findQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)

	s.Run("should update a receiver", func() {
		timeNow := time.Now()
		expectedReceiver := &domain.Receiver{
			Id:             10,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		input := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		expectedRows2 := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow("baz", expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Id, "baz", input.Type,
			json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(findQuery).WillReturnRows(expectedRows2)

		err := s.repository.Update(input)
		s.Equal("baz", input.Name)
		s.Nil(err)
	})

	s.Run("should return error if receiver does not exist", func() {
		timeNow := time.Now()
		input := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		s.dbmock.ExpectExec(updateQuery).WithArgs(input.Id, "baz", input.Type,
			json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repository.Update(input)
		s.EqualError(err, "receiver doesn't exist")
	})

	s.Run("should return error updating the receiver", func() {
		timeNow := time.Now()
		expectedReceiver := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,

			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		input := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		s.dbmock.ExpectExec(updateQuery).WithArgs(expectedReceiver.Id, "baz", expectedReceiver.Type,
			json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			AnyTime{}, AnyTime{}, expectedReceiver.Id, expectedReceiver.Id).
			WillReturnError(errors.New("random error"))

		err := s.repository.Update(input)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated receiver", func() {
		timeNow := time.Now()
		expectedReceiver := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		input := &domain.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		s.dbmock.ExpectExec(updateQuery).WithArgs(expectedReceiver.Id, "baz", expectedReceiver.Type,
			json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`),
			AnyTime{}, AnyTime{}, expectedReceiver.Id, expectedReceiver.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(findQuery).WillReturnError(errors.New("random error"))

		err := s.repository.Update(input)
		s.EqualError(err, "random error")
	})
}

func (s *ReceiverRepositoryTestSuite) TestDelete() {
	s.Run("should delete receiver of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "receivers" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.repository.Delete(1)
		s.Nil(err)
	})

	s.Run("should return error in deleting receiver of given id", func() {
		expectedQuery := regexp.QuoteMeta(`DELETE FROM "receivers" WHERE id = $1`)
		s.dbmock.ExpectExec(expectedQuery).WillReturnError(errors.New("random error"))

		err := s.repository.Delete(1)
		s.EqualError(err, "random error")
	})
}

func TestReceiverRepository(t *testing.T) {
	suite.Run(t, new(ReceiverRepositoryTestSuite))
}
