package notification

import (
	"context"
	"sync"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/errors"
)

type Worker struct {
	logger   log.Logger
	handlers []*Handler
}

func NewWorker(logger log.Logger, handlers ...*Handler) *Worker {
	return &Worker{
		logger:   logger,
		handlers: handlers,
	}
}

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

	w.logger.Info("all workers exited")
	err := ctx.Err()
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
