package postgres_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store"
	"github.com/odpf/siren/internal/store/postgres"
	"github.com/odpf/siren/mocks"
	"github.com/stretchr/testify/suite"
)

type SubscriptionRepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository store.SubscriptionRepository
}

func (s *SubscriptionRepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = postgres.NewSubscriptionRepository(db)
}

func (s *SubscriptionRepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *SubscriptionRepositoryTestSuite) TestCreate() {
	match := make(map[string]string)
	inputRandomConfig := make(map[string]string)
	randomSlackReceiverConfig := make(map[string]interface{})
	randomPagerdutyReceiverConfig := make(map[string]interface{})
	randomHTTPReceiverConfig := make(map[string]interface{})
	match["foo"] = "bar"
	inputRandomConfig["channel_name"] = "test"
	randomSlackReceiverConfig["token"] = "xoxb"
	randomPagerdutyReceiverConfig["service_key"] = "abcd"
	randomHTTPReceiverConfig["url"] = "http://localhost:3000"
	receiver1 := domain.ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
	receiver2 := domain.ReceiverMetadata{Id: 2, Configuration: make(map[string]string)}
	input := &domain.Subscription{
		Id:        1,
		Namespace: 1,
		Urn:       "foo",
		Match:     match,
		Receivers: []domain.ReceiverMetadata{receiver2, receiver1},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	rawReceivers := `[{"id":2,"configuration":{}},{"id":1,"configuration":{"channel_name":"test"}}]`

	insertQuery := regexp.QuoteMeta(`INSERT INTO "subscriptions" ("namespace_id","urn","receiver","match","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)

	s.Run("should create a subscription", func() {
		expectedID := uint64(1)
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(input.Namespace, input.Urn,
				rawReceivers, json.RawMessage(`{"foo":"bar"}`), input.CreatedAt,
				input.UpdatedAt, input.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

		err := s.repository.Create(context.Background(), input)
		s.Nil(err)
		s.Equal(expectedID, input.Id)
	})

	s.Run("should create a subscription with transaction", func() {
		expectedID := uint64(1)
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(input.Namespace, input.Urn,
				rawReceivers, json.RawMessage(`{"foo":"bar"}`), input.CreatedAt,
				input.UpdatedAt, input.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
		s.dbmock.ExpectCommit()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Create(ctx, input)
		commitErr := s.repository.Commit(ctx)
		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(expectedID, input.Id)
	})

	s.Run("should rollback transaction", func() {
		expectedID := uint64(1)
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(input.Namespace, input.Urn,
				rawReceivers, json.RawMessage(`{"foo":"bar"}`), input.CreatedAt,
				input.UpdatedAt, input.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
		s.dbmock.ExpectRollback()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Create(ctx, input)
		commitErr := s.repository.Rollback(ctx)
		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(expectedID, input.Id)
	})

	s.Run("should return error from creation", func() {
		s.dbmock.ExpectQuery(insertQuery).
			WithArgs(input.Namespace, input.Urn,
				rawReceivers, json.RawMessage(`{"foo":"bar"}`), input.CreatedAt,
				input.UpdatedAt, input.Id).
			WillReturnError(errors.New("random error"))

		err := s.repository.Create(context.Background(), input)
		s.EqualError(err, "failed to insert subscription: random error")
	})
}

func (s *SubscriptionRepositoryTestSuite) TestGet() {
	expectedSubscription := &domain.Subscription{
		Id:        1,
		Namespace: 1,
		Urn:       "foo",
		Match:     make(map[string]string),
		Receivers: []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.Run("should get subscription by id", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)

		actualSubscription, err := s.repository.Get(context.Background(), 1)

		s.Equal(expectedSubscription, actualSubscription)
		s.Nil(err)
	})

	s.Run("should get subscription by id using transcation", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)

		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectCommit()

		ctx := s.repository.WithTransaction(context.Background())
		actualSubscription, err := s.repository.Get(ctx, 1)
		commitErr := s.repository.Commit(ctx)

		s.Nil(commitErr)
		s.Equal(expectedSubscription, actualSubscription)
		s.Nil(err)
	})

	s.Run("should rollback transcation", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)

		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectRollback()

		ctx := s.repository.WithTransaction(context.Background())
		actualSubscription, err := s.repository.Get(ctx, 1)
		commitErr := s.repository.Rollback(ctx)

		s.Nil(commitErr)
		s.Equal(expectedSubscription, actualSubscription)
		s.Nil(err)
	})

	s.Run("should return error from db", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))

		actualSubscription, err := s.repository.Get(context.Background(), 1)
		s.Nil(actualSubscription)
		s.EqualError(err, "random error")
	})

	s.Run("should return nil if subscription not found", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"})
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)

		actualSubscription, err := s.repository.Get(context.Background(), 1)
		s.Nil(actualSubscription)
		s.Nil(err)
	})
}

func (s *SubscriptionRepositoryTestSuite) TestList() {
	expectedSubscription := &domain.Subscription{
		Id:        1,
		Namespace: 1,
		Urn:       "foo",
		Match:     make(map[string]string),
		Receivers: []domain.ReceiverMetadata{{Id: 1, Configuration: make(map[string]string)}},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.Run("should get all subscriptions", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualSubscriptions, err := s.repository.List(context.Background())

		s.Equal(1, len(actualSubscriptions))
		s.Equal(expectedSubscription.Id, actualSubscriptions[0].Id)
		s.Equal(expectedSubscription.Namespace, actualSubscriptions[0].Namespace)
		s.Equal(expectedSubscription.Urn, actualSubscriptions[0].Urn)
		s.Equal(expectedSubscription.Match, actualSubscriptions[0].Match)
		s.Nil(err)
	})

	s.Run("should get all subscriptions using transaction", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectCommit()

		ctx := s.repository.WithTransaction(context.Background())
		actualSubscriptions, err := s.repository.List(ctx)
		commitErr := s.repository.Commit(ctx)

		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(1, len(actualSubscriptions))
		s.Equal(expectedSubscription, actualSubscriptions[0])
	})

	s.Run("should rollback transaction", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.Namespace,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectRollback()

		ctx := s.repository.WithTransaction(context.Background())
		actualSubscriptions, err := s.repository.List(ctx)
		commitErr := s.repository.Rollback(ctx)

		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(1, len(actualSubscriptions))
		s.Equal(expectedSubscription, actualSubscriptions[0])
	})

	s.Run("should return error in fetching subscriptions", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		actualSubscriptions, err := s.repository.List(context.Background())

		s.EqualError(err, "random error")
		s.Nil(actualSubscriptions)
	})
}

func (s *SubscriptionRepositoryTestSuite) TestUpdate() {
	timeNow := time.Now()
	match := make(map[string]string)
	inputRandomConfig := make(map[string]string)
	randomSlackReceiverConfig := make(map[string]interface{})
	randomPagerdutyReceiverConfig := make(map[string]interface{})
	randomHTTPReceiverConfig := make(map[string]interface{})
	match["foo"] = "bar"
	inputRandomConfig["channel_name"] = "test"
	randomSlackReceiverConfig["token"] = "xoxb"
	randomPagerdutyReceiverConfig["service_key"] = "abcd"
	randomHTTPReceiverConfig["url"] = "http://localhost:3000"
	receiver := domain.ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
	subscription := &domain.Subscription{
		Id:        1,
		Namespace: 1,
		Urn:       "foo",
		Match:     match,
		Receivers: []domain.ReceiverMetadata{receiver},
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	inputRandomConfig["channel_name"] = "updated_channel"

	input := &domain.Subscription{
		Id:        1,
		Namespace: 1,
		Urn:       "foo",
		Match:     match,
		Receivers: []domain.ReceiverMetadata{receiver},
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	rawReceivers := `[{"id":1,"configuration":{"channel_name":"updated_channel"}}]`
	updateQuery := regexp.QuoteMeta(`UPDATE "subscriptions" SET "namespace_id"=$1,"urn"=$2,"receiver"=$3,"match"=$4,"created_at"=$5,"updated_at"=$6 WHERE id = $7 AND "id" = $8`)
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1 AND "subscriptions"."id" = $1`)

	s.Run("should update a subscription", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(input.Namespace, input.Urn, rawReceivers, json.RawMessage(`{"foo":"bar"}`),
				AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{"channel_name":"updated_channel"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.Id).
			WillReturnRows(expectedRows)

		err := s.repository.Update(context.Background(), input)
		s.Nil(err)
		s.Equal(subscription.Receivers[0].Configuration, input.Receivers[0].Configuration)
	})

	s.Run("should update a subscription with transaction", func() {
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(input.Namespace, input.Urn, rawReceivers, json.RawMessage(`{"foo":"bar"}`),
				AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{"channel_name":"updated_channel"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.Id).
			WillReturnRows(expectedRows)
		s.dbmock.ExpectCommit()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Update(ctx, input)
		commitErr := s.repository.Commit(ctx)

		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(subscription.Receivers[0].Configuration, input.Receivers[0].Configuration)
	})

	s.Run("should rollback transaction", func() {
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(input.Namespace, input.Urn, rawReceivers, json.RawMessage(`{"foo":"bar"}`),
				AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.Namespace,
				json.RawMessage(`[{"id":1,"configuration":{"channel_name":"updated_channel"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.Id).
			WillReturnRows(expectedRows)
		s.dbmock.ExpectRollback()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Update(ctx, input)
		commitErr := s.repository.Rollback(ctx)

		s.Nil(commitErr)
		s.Nil(err)
		s.Equal(subscription.Receivers[0].Configuration, input.Receivers[0].Configuration)
	})

	s.Run("should return error if subscription does not exist", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(input.Namespace, input.Urn, rawReceivers, json.RawMessage(`{"foo":"bar"}`),
				AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.repository.Update(context.Background(), input)
		s.EqualError(err, "subscription doesn't exist")
	})

	s.Run("should return error if got error fetching updated subscription", func() {
		s.dbmock.ExpectExec(updateQuery).
			WithArgs(input.Namespace, input.Urn, rawReceivers, json.RawMessage(`{"foo":"bar"}`),
				AnyTime{}, AnyTime{}, input.Id, input.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.dbmock.ExpectQuery(selectQuery).WithArgs(input.Id).
			WillReturnError(errors.New("random error"))

		err := s.repository.Update(context.Background(), input)
		s.EqualError(err, "random error")
	})
}

func (s *SubscriptionRepositoryTestSuite) TestDelete() {
	deleteQuery := regexp.QuoteMeta(`DELETE FROM "subscriptions" WHERE "subscriptions"."id" = $1`)

	s.Run("should delete a subscription", func() {
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
		err := s.repository.Delete(context.Background(), 1)
		s.Nil(err)
	})

	s.Run("should delete a subscription with transaction", func() {
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
		s.dbmock.ExpectCommit()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Delete(ctx, 1)
		commitErr := s.repository.Commit(ctx)
		s.Nil(commitErr)
		s.Nil(err)
	})

	s.Run("should rollback transaction", func() {
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
		s.dbmock.ExpectRollback()

		ctx := s.repository.WithTransaction(context.Background())
		err := s.repository.Delete(ctx, 1)
		commitErr := s.repository.Rollback(ctx)
		s.Nil(commitErr)
		s.Nil(err)
	})

	s.Run("should return error from db", func() {
		s.dbmock.ExpectExec(deleteQuery).WillReturnError(errors.New("random error"))
		err := s.repository.Delete(context.Background(), 1)
		s.EqualError(err, "random error")
	})
}

func TestSubscriptionRepository(t *testing.T) {
	suite.Run(t, new(SubscriptionRepositoryTestSuite))
}
