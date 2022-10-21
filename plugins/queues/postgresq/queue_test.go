package postgresq_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertest"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/plugins/queues/postgresq"
	"github.com/odpf/siren/plugins/queues/postgresq/migrations"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	logger   log.Logger
	ctx      context.Context
	dbc      *db.Client
	pool     *dockertest.Pool
	resource *dockertest.Resource
	q        *postgresq.Queue
	dlq      *postgresq.Queue
}

func (s *QueueTestSuite) SetupSuite() {
	var (
		err      error
		pgUser   = "test_user"
		pgPass   = "test_pass"
		pgDBName = "test_db"
	)

	s.logger = log.NewZap()
	dpg, err := dockertest.CreatePostgres(
		dockertest.PostgresWithDetail(
			pgUser, pgPass, pgDBName,
		),
	)
	if err != nil {
		s.T().Fatal(err)
	}

	s.pool = dpg.GetPool()
	s.resource = dpg.GetResource()

	dbConfig := db.Config{
		Driver: "postgres",
	}
	dbConfig.URL = dpg.GetExternalConnString()
	s.dbc, err = db.New(dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}

	s.ctx = context.TODO()
	err = db.RunMigrations(dbConfig, migrations.FS, migrations.ResourcePath)
	if err != nil {
		s.T().Fatal(err)
	}

	s.q, err = postgresq.New(s.logger, dbConfig)
	if err != nil {
		s.T().Fatal(err)
	}

	s.dlq, err = postgresq.New(s.logger, dbConfig, postgresq.WithStrategy(postgresq.StrategyDLQ))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *QueueTestSuite) TearDownSuite() {
	s.q.Stop(s.ctx)
	// Clean tests
	if err := s.pool.Purge(s.resource); err != nil {
		s.T().Fatal(err)
	}
}

func (s *QueueTestSuite) cleanup() error {
	_, err := s.dbc.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", postgresq.MESSAGE_QUEUE_TABLE_NAME))
	if err != nil {
		return err
	}
	return nil
}

func (s *QueueTestSuite) TestSimpleEnqueueDequeue() {
	timeNow := time.Now()
	ns := []notification.Notification{
		{
			ID:        uuid.NewString(),
			Data:      map[string]interface{}{},
			Labels:    map[string]string{},
			CreatedAt: timeNow,
		},
		{
			ID:        uuid.NewString(),
			Data:      map[string]interface{}{},
			Labels:    map[string]string{},
			CreatedAt: timeNow,
		},
		{
			ID:        uuid.NewString(),
			Data:      map[string]interface{}{},
			Labels:    map[string]string{},
			CreatedAt: timeNow,
		},
		{
			ID:        uuid.NewString(),
			Data:      map[string]interface{}{},
			Labels:    map[string]string{},
			CreatedAt: timeNow,
		},
	}

	messages := []notification.Message{}
	for _, n := range ns {
		msg, err := n.ToMessage(receiver.TypeSlack, map[string]interface{}{})
		s.Require().NoError(err)
		messages = append(messages, *msg)
	}

	s.Run("should return no error if all messages are successfully processed", func() {
		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			s.Assert().Len(messages, 1)
			return nil
		}

		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		for i := 0; i < len(messages); i++ {
			_ = s.q.Dequeue(s.ctx, nil, 1, handlerFn)
		}

		err = s.cleanup()
		s.Require().NoError(err)
	})

	s.Run("should return no error if all messages are successfully processed with different batch", func() {
		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			s.Assert().Len(messages, 2)
			return nil
		}

		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		for i := 0; i < 2; i++ {
			_ = s.q.Dequeue(s.ctx, nil, 2, handlerFn)
		}

		err = s.cleanup()
		s.Require().NoError(err)
	})

	s.Run("should return an error if a message is failed to process", func() {
		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			return errors.New("some error")
		}

		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		for i := 0; i < len(messages); i++ {
			err := s.q.Dequeue(s.ctx, nil, 1, handlerFn)
			s.Assert().Error(errors.New("error processing dequeued message: some error"), err)
		}

		err = s.cleanup()
		s.Require().NoError(err)
	})
}

func (s *QueueTestSuite) TestEnqueueDequeueWithCallback() {
	messages := make([]notification.Message, 5)

	for i := 0; i < len(messages); i++ {
		messages[i].Initialize(notification.Notification{}, receiver.TypeSlack, map[string]interface{}{}, notification.InitWithID(fmt.Sprintf("%d", i+1)))
	}

	s.Run("should update row with error for id \"5\"", func() {
		var anError = errors.New("some error")

		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		for _, m := range messages {
			if m.ID == "5" {
				m.MarkFailed(time.Now(), true, anError)
				err = s.q.ErrorCallback(s.ctx, m)
				s.Assert().NoError(err)
			}
		}

		tempMessage := &postgresq.NotificationMessage{}
		err = s.dbc.Get(tempMessage, fmt.Sprintf("SELECT * FROM %s WHERE id = '5'", postgresq.MESSAGE_QUEUE_TABLE_NAME))
		s.Require().NoError(err)

		s.Assert().Equal(string(notification.MessageStatusFailed), tempMessage.Status)
		s.Assert().Equal(anError.Error(), tempMessage.LastError.String)
		s.Assert().Equal(1, tempMessage.TryCount)

		err = s.cleanup()
		s.Require().NoError(err)
	})

	s.Run("should update row with when successfully published", func() {
		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		for _, m := range messages {
			m.MarkPublished(time.Now())
			err = s.q.SuccessCallback(s.ctx, m)
			s.Assert().NoError(err)
		}

		tempMessage := &postgresq.NotificationMessage{}
		err = s.dbc.Get(tempMessage, fmt.Sprintf("SELECT * FROM %s LIMIT 1", postgresq.MESSAGE_QUEUE_TABLE_NAME))
		s.Require().NoError(err)

		s.Assert().Equal(string(notification.MessageStatusPublished), tempMessage.Status)
		s.Assert().Equal(1, tempMessage.TryCount)

		err = s.cleanup()
		s.Require().NoError(err)
	})
}

func (s *QueueTestSuite) TestEnqueueDequeueDLQ() {
	messages := make([]notification.Message, 5)

	for i := 0; i < len(messages); i++ {
		messages[i].Initialize(notification.Notification{}, receiver.TypeSlack, map[string]interface{}{}, notification.InitWithID(fmt.Sprintf("%d", i+1)))
	}

	s.Run("failed messages should be re-processed by dlq an ignored by main queue", func() {
		var anError = errors.New("some error")

		err := s.q.Enqueue(s.ctx, messages...)
		s.Require().NoError(err)

		// mark failed all
		for _, m := range messages {
			m.MarkFailed(time.Now(), true, anError)
			err = s.q.ErrorCallback(s.ctx, m)
			s.Assert().NoError(err)
		}

		_ = s.q.Dequeue(s.ctx, nil, 5, func(ctx context.Context, m []notification.Message) error { s.Assert().Empty(m); return nil })
		s.Assert().NoError(err)

		_ = s.dlq.Dequeue(s.ctx, nil, 5, func(ctx context.Context, m []notification.Message) error {
			s.Assert().Len(m, 5)
			return nil
		})

		tempMessage := &postgresq.NotificationMessage{}
		err = s.dbc.Get(tempMessage, fmt.Sprintf("SELECT * FROM %s LIMIT 1", postgresq.MESSAGE_QUEUE_TABLE_NAME))
		s.Require().NoError(err)

		s.Assert().Equal(string(notification.MessageStatusPending), tempMessage.Status)
		s.Assert().Equal(anError.Error(), tempMessage.LastError.String)
		s.Assert().Equal(1, tempMessage.TryCount)

		err = s.cleanup()
		s.Require().NoError(err)
	})
}

func TestQueue(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
