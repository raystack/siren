package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/goto/salt/log"
	"github.com/goto/siren/core/notification"
	"github.com/goto/siren/plugins"
	"github.com/goto/siren/plugins/queues"
)

// Queue simulates queue inmemory, this is for testing only
// not recommended to use this in production
type Queue struct {
	logger     log.Logger
	once       sync.Once
	stopSignal chan struct{}
	memoryQ    chan notification.Message
}

// New creates a new queue instance
func New(logger log.Logger, capacity uint) *Queue {
	return &Queue{
		logger:     logger,
		stopSignal: make(chan struct{}),
		memoryQ:    make(chan notification.Message, capacity),
	}
}

// Dequeue pop the queue based on specific filters (receiver types or batch size)
// and process the messages with handlerFn
func (q *Queue) Dequeue(ctx context.Context, receiverTypes []string, batchSize int, handlerFn func(context.Context, []notification.Message) error) error {
	messages := []notification.Message{}
	for i := 0; i < batchSize; i++ {
		var message notification.Message
		select {
		case <-ctx.Done():
			q.logger.Info("inmemory dequeue work is done", "scope", "queues.inmemory.dequeue")
			return nil
		case message = <-q.memoryQ:
			q.logger.Debug("dequeued a message")
		default:
			q.logger.Debug("queue empty")
			return notification.ErrNoMessage
		}

		messages = append(messages, message)
	}

	if err := handlerFn(ctx, messages); err != nil {
		return fmt.Errorf("error processing dequeued message: %w", err)
	}

	return nil
}

// Enqueue pushes messages to the queue
func (q *Queue) Enqueue(ctx context.Context, ms ...notification.Message) error {
	for _, m := range ms {
		select {
		case <-q.stopSignal:
			q.logger.Debug("enqueuer retrieving stop signal")
			return nil
		case q.memoryQ <- m:
			q.logger.Debug("enqueued message", "scope", "queues.inmemory.enqueue", "type", m.ReceiverType, "configs", m.Configs, "details", m.Details)
			continue
		default:
			return fmt.Errorf("error enqueueing message: %v", m.Details)
		}
	}
	return nil
}

// SuccessCallback is a callback that will be called once the message is succesfully handled by handlerFn
func (q *Queue) SuccessCallback(ctx context.Context, ms notification.Message) error {
	q.logger.Debug("successfully sending message", "scope", "queues.inmemory.success_callback", "type", ms.ReceiverType, "configs", ms.Configs, "details", ms.Details)
	return nil
}

// ErrorCallback is a callback that will be called once the message is failed to be handled by handlerFn
func (q *Queue) ErrorCallback(ctx context.Context, ms notification.Message) error {
	q.logger.Error("failed sending message", "scope", "queues.inmemory.error_callback", "type", ms.ReceiverType, "configs", ms.Configs, "details", ms.Details, "last_error", ms.LastError)
	return nil
}

func (q *Queue) Cleanup(ctx context.Context, filter queues.FilterCleanup) error {
	return plugins.ErrNotImplemented
}

func (q *Queue) Type() string {
	return "inmemory"
}

// Stop is a inmemmory queue function
// this will close the channel to simulate queue
func (q *Queue) Stop(ctx context.Context) error {
	q.once.Do(func() {
		q.logger.Debug("closing inmemory queue channel")
		close(q.memoryQ)
	})
	return nil
}
