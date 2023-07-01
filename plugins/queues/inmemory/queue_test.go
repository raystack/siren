package inmemory_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/raystack/salt/log"
	"github.com/raystack/siren/core/notification"
	"github.com/raystack/siren/plugins/queues/inmemory"
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
		q := inmemory.New(logger, 10)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			assert.Len(t, messages, 1)
			return nil
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < len(messages); i++ {
			_ = q.Dequeue(ctx, nil, 1, handlerFn)
		}

		wg.Wait()

		q.Stop(ctx)
	})

	t.Run("should return no error if all messages are successfully processed with different batch", func(t *testing.T) {
		ctx := context.Background()
		logger := log.NewZap()
		q := inmemory.New(logger, 10)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			assert.Len(t, messages, 2)
			return nil
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < 2; i++ {
			_ = q.Dequeue(ctx, nil, 2, handlerFn)
		}

		wg.Wait()

		q.Stop(ctx)
	})

	t.Run("should return an error if a message is failed to process", func(t *testing.T) {
		ctx := context.Background()
		logger := log.NewZap()
		q := inmemory.New(logger, 10)

		handlerFn := func(ctx context.Context, messages []notification.Message) error {
			return errors.New("some error")
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := q.Enqueue(ctx, messages...)
			require.NoError(t, err)
		}()

		for i := 0; i < len(messages); i++ {
			err := q.Dequeue(ctx, nil, 1, handlerFn)
			assert.Error(t, errors.New("error processing dequeued message: some error"), err)
		}

		wg.Wait()

		q.Stop(ctx)
	})
}
