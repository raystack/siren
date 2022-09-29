package inmemory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/notification"
	"github.com/odpf/siren/plugins/queues/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	messages := []notification.Message{
		{
			ID: "1",
		},
		{
			ID: "2",
		},
		{
			ID: "3",
		},
		{
			ID: "4",
		},
	}

	t.Run("should return no error if all messages are successfully processed", func(t *testing.T) {
		ctx := context.Background()
		logger := log.NewZap()
		q := inmemory.New(logger)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			assert.Len(t, messages, 1)
			return nil
		}

		go func() {
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < len(messages); i++ {
			err := q.Dequeue(ctx, nil, 1, handlerFn)
			require.NoError(t, err)
		}

		q.Close()
	})

	t.Run("should return no error if all messages are successfully processed with different batch", func(t *testing.T) {
		ctx := context.Background()
		logger := log.NewZap()
		q := inmemory.New(logger)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			assert.Len(t, messages, 2)
			return nil
		}

		go func() {
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < 2; i++ {
			err := q.Dequeue(ctx, nil, 2, handlerFn)
			require.NoError(t, err)
		}

		q.Close()
	})

	t.Run("should return an error if a message is failed to process", func(t *testing.T) {
		ctx := context.Background()
		logger := log.NewZap()
		q := inmemory.New(logger)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			return errors.New("some error")
		}

		go func() {
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < len(messages); i++ {
			err := q.Dequeue(ctx, nil, 1, handlerFn)
			assert.Error(t, errors.New("error processing dequeued message: some error"), err)
		}

		q.Close()
	})
}
