package worker

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/goto/salt/log"
)

const (
	defaultPollDuration = 5 * time.Second
)

// Ticker is a worker that runs periodically
type Ticker struct {
	id           string
	logger       log.Logger
	pollDuration time.Duration
}

// NewTicker creates a new worker that does an action periodically
func NewTicker(logger log.Logger, opts ...TickerOption) *Ticker {
	wt := &Ticker{
		logger: logger,
	}

	for _, opt := range opts {
		opt(wt)
	}

	if wt.id == "" {
		wt.id = uuid.NewString()
	}

	if wt.pollDuration == 0 {
		wt.pollDuration = defaultPollDuration
	}

	return wt
}

// Run starts worker that handle a task periodically
func (wt *Ticker) Run(ctx context.Context, cancelChan chan struct{}, handlerFn func(ctx context.Context, runningAt time.Time) error) {
	ticker := time.NewTicker(wt.pollDuration)
	defer ticker.Stop()

	wt.logger.Info("running worker", "id", wt.id)

	for {
		select {
		case <-cancelChan:
			wt.logger.Info("stopping worker", "id", wt.id)
			return

		case t := <-ticker.C:
			if err := handlerFn(ctx, t); err != nil {
				wt.logger.Error("error running worker", "error", err, "id", wt.id)
			}
		}
	}
}

// GetID fetch identifier of a worker
func (wt *Ticker) GetID() string {
	return wt.id
}
