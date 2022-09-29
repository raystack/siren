package inmemory

import (
	"context"
	"fmt"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/notification"
)

// Queue simulates queue inmemory, this is for testing only
// not recommended to use this in production
type Queue struct {
	logger  log.Logger
	memoryQ chan notification.Message
}

// New creates a new queue instance
func New(logger log.Logger) *Queue {
	return &Queue{
		logger:  logger,
		memoryQ: make(chan notification.Message),
	}
}

// Dequeue pop the queue based on specific filters (receiver types or batch size)
// and process the messages with handlerFn
func (q *Queue) Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []notification.Message) error) error {
	cancelableCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	messages := []notification.Message{}
	for i := 0; i < batchSize; i++ {
		select {
		case <-cancelableCtx.Done():
			q.logger.Info("inmemory dequeue work is done", "scope", "queues.inmemory.dequeue")
			return nil

		case message := <-q.memoryQ:
			messages = append(messages, message)
		}
	}

	if err := handlerFn(cancelableCtx, messages); err != nil {
		return fmt.Errorf("error processing dequeued message: %w", err)
	}

	return nil
}

// Enqueue pushes messages to the queue
func (q *Queue) Enqueue(ctx context.Context, ms ...notification.Message) error {
	for _, m := range ms {
		q.memoryQ <- m
		q.logger.Debug("enqueued message", "scope", "queues.inmemory.enqueue", "type", m.ReceiverType, "configs", m.Configs, "detail", m.Detail)
	}

	return nil
}

// SuccessHandler is a callback that will be called once the message is succesfully handled by handlerFn
func (q *Queue) SuccessHandler(ctx context.Context, ms notification.Message) error {
	q.logger.Debug("successfully sending message", "scope", "queues.inmemory.success_handler", "type", ms.ReceiverType, "configs", ms.Configs, "detail", ms.Detail)
	return nil
}

// ErrorHandler is a callback that will be called once the message is failed to be handled by handlerFn
func (q *Queue) ErrorHandler(ctx context.Context, ms notification.Message) error {
	q.logger.Error("failed sending message", "scope", "queues.inmemory.error_handler", "type", ms.ReceiverType, "configs", ms.Configs, "detail", ms.Detail, "last_error", ms.LastError)
	return nil
}

// Close is a specific inmemmory queue function
// this will close the channel to simulate queue
func (q *Queue) Close() {
	close(q.memoryQ)
}
