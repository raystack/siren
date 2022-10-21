package retry

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Runner knows how to execute a execution logic and returns error if errors.
type Runner interface {
	// Run will run the unit of execution passed on f.
	Run(ctx context.Context, f func(ctx context.Context) error) error
}

type wrapper struct {
	c Config
}

func New(cfg Config) Runner {
	if cfg.WaitDuration <= 0 {
		cfg.WaitDuration = 20 * time.Millisecond
	}

	if cfg.MaxTries <= 0 {
		cfg.MaxTries = 3
	}

	return &wrapper{
		c: cfg,
	}
}

// Run wraps function call with retry logic
// exponential backoff option is configurable
func (w *wrapper) Run(ctx context.Context, f func(ctx context.Context) error) error {
	if !w.c.Enable {
		return f(ctx)
	}

	var err error

	for i := 0; i <= w.c.MaxTries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
			err = f(ctx)
			if err == nil || !errors.As(err, new(RetryableError)) {
				return err
			}

			waitDuration := w.c.WaitDuration

			// Apply Backoff.
			// The backoff is calculated exponentially based on a base time
			// and the attemp of the retry.
			if w.c.EnableBackoff {
				exp := math.Exp2(float64(i + 1))
				waitDuration = time.Duration(float64(w.c.WaitDuration) * exp)

				// Round to millisecs.
				waitDuration = waitDuration.Round(time.Millisecond)

				// "full jitter".
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				waitDuration = time.Duration(float64(waitDuration) * random.Float64())
			}
			time.Sleep(waitDuration)
		}
	}
	return err
}
