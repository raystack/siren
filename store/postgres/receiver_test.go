package postgres_test

import (
	"database/sql"
	"encoding/json"
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
		configurations := make(model.StringInterfaceMap)
		configurations["foo"] = "bar"
		labels := make(model.StringStringMap)
		labels["foo"] = "bar"

		receiver := &model.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		expectedReceivers := []*model.Receiver{receiver}

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
	configurations := make(model.StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should create a receiver", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "receivers"
											("name","type","labels","configurations","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)
		expectedReceiver := &model.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedReceiver.Name,
			expectedReceiver.Type, expectedReceiver.Labels, expectedReceiver.Configurations,
			expectedReceiver.CreatedAt, expectedReceiver.UpdatedAt, expectedReceiver.Id).
			WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow(expectedReceiver.Name, expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)

		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualReceiver, err := s.repository.Create(expectedReceiver)
		s.Equal(expectedReceiver, actualReceiver)
		s.Nil(err)
	})

	s.Run("should return errors in creating a receiver", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "receivers"
											("name","type","labels","configurations","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		expectedReceiver := &model.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedReceiver.Name, expectedReceiver.Type,
			expectedReceiver.Labels, expectedReceiver.Configurations,
			expectedReceiver.CreatedAt, expectedReceiver.UpdatedAt, expectedReceiver.Id).
			WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Create(expectedReceiver)
		s.EqualError(err, "random error")
		s.Nil(actualReceiver)
	})

	s.Run("should return error if finding newly inserted receiver fails", func() {
		insertQuery := regexp.QuoteMeta(`INSERT INTO "receivers"
											("name","type","labels","configurations","created_at","updated_at","id")
											VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)
		expectedReceiver := &model.Receiver{
			Id:             1,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedReceiver.Name,
			expectedReceiver.Type, expectedReceiver.Labels, expectedReceiver.Configurations,
			expectedReceiver.CreatedAt, expectedReceiver.UpdatedAt, expectedReceiver.Id).
			WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Create(expectedReceiver)
		s.EqualError(err, "random error")
		s.Nil(actualReceiver)
	})
}

func (s *ReceiverRepositoryTestSuite) TestGet() {
	configurations := make(model.StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should get receiver by id", func() {
		expectedQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 1`)
		expectedReceiver := &model.Receiver{
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
	configurations := make(model.StringInterfaceMap)
	configurations["foo"] = "bar"
	labels := make(model.StringStringMap)
	labels["foo"] = "bar"

	s.Run("should update a receiver", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "receivers"
						SET "name"=$1,"type"=$2,"labels"=$3,"configurations"=$4,"created_at"=$5,"updated_at"=$6
						WHERE id = $7 AND "id" = $8`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		timeNow := time.Now()
		expectedReceiver := &model.Receiver{
			Id:             10,
			Name:           "foo",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		input := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		expectedRows1 := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow(expectedReceiver.Name, expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		expectedRows2 := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow("baz", expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows1)
		s.dbmock.ExpectExec(updateQuery).WithArgs("baz", expectedReceiver.Type,
			expectedReceiver.Labels, expectedReceiver.Configurations,
			AnyTime{}, AnyTime{}, expectedReceiver.Id, expectedReceiver.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRows2)

		actualReceiver, err := s.repository.Update(input)
		s.Equal("baz", actualReceiver.Name)
		s.Nil(err)
	})

	s.Run("should return error if receiver does not exist", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		timeNow := time.Now()
		input := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.NewRows(nil))

		actualReceiver, err := s.repository.Update(input)
		s.Nil(actualReceiver)
		s.EqualError(err, "receiver doesn't exist")
	})

	s.Run("should return error in finding the receiver", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		timeNow := time.Now()
		input := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Update(input)
		s.Nil(actualReceiver)
		s.EqualError(err, "random error")
	})

	s.Run("should return error updating the receiver", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "receivers"
						SET "name"=$1,"type"=$2,"labels"=$3,"configurations"=$4,"created_at"=$5,"updated_at"=$6
						WHERE id = $7 AND "id" = $8`)
		timeNow := time.Now()
		expectedReceiver := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,

			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		input := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		expectedRows := sqlmock.
			NewRows([]string{"urn", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow(expectedReceiver.Name, expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs("baz", expectedReceiver.Type, expectedReceiver.Labels,
			expectedReceiver.Configurations, AnyTime{}, AnyTime{}, expectedReceiver.Id, expectedReceiver.Id).
			WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Update(input)
		s.Nil(actualReceiver)
		s.EqualError(err, "random error")
	})

	s.Run("should return error in finding the updated receiver", func() {
		firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		updateQuery := regexp.QuoteMeta(`UPDATE "receivers"
						SET "name"=$1,"type"=$2,"labels"=$3,"configurations"=$4,"created_at"=$5,"updated_at"=$6
						WHERE id = $7 AND "id" = $8`)
		secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "receivers" WHERE id = 10`)
		timeNow := time.Now()
		expectedReceiver := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}
		input := &model.Receiver{
			Id:             10,
			Name:           "baz",
			Type:           "slack",
			Labels:         labels,
			Configurations: configurations,
			CreatedAt:      timeNow,
			UpdatedAt:      timeNow,
		}

		expectedRows := sqlmock.
			NewRows([]string{"name", "type", "labels", "configurations", "created_at", "updated_at", "id"}).
			AddRow(expectedReceiver.Name, expectedReceiver.Type,
				json.RawMessage(`{"foo": "bar"}`), json.RawMessage(`{"foo": "bar"}`), expectedReceiver.CreatedAt,
				expectedReceiver.UpdatedAt, expectedReceiver.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(updateQuery).WithArgs("baz",
			expectedReceiver.Type, expectedReceiver.Labels, expectedReceiver.Configurations,
			AnyTime{}, AnyTime{}, expectedReceiver.Id, expectedReceiver.Id).
			WillReturnResult(sqlmock.NewResult(10, 1))
		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))

		actualReceiver, err := s.repository.Update(input)
		s.Nil(actualReceiver)
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
