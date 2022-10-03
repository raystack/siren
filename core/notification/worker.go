package notification

import (
	"context"
	"sync"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/errors"
)

// Worker is a notification worker instance that runs one or more than one
// notification handler
type Worker struct {
	logger   log.Logger
	handlers []*Handler
}

// NewWorker creates a new worker instance
func NewWorker(logger log.Logger, handlers ...*Handler) *Worker {
	return &Worker{
		logger:   logger,
		handlers: handlers,
	}
}

// Run will execute and run one or multiple notification handlers
// as goroutines
func (w *Worker) Run(ctx context.Context) error {
	cancellableCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}
	for _, handler := range w.handlers {
		wg.Add(1)
		go func(h *Handler) {
			defer wg.Done()
			w.logger.Info("running handler worker", "id", h.id)
			h.RunHandler(cancellableCtx)
			w.logger.Info("handler worker exited", "id", h.id)
		}(handler)
	}
	wg.Wait()

	w.logger.Info("all handlers exited")
	err := ctx.Err()
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
