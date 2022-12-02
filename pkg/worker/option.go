package worker

import "time"

// TickerOption is an option to customize worker ticker creation
type TickerOption func(*Ticker)

// WithTickerDuration sets created handler with the specified poll duration
func WithTickerDuration(pollDuration time.Duration) TickerOption {
	return func(wt *Ticker) {
		wt.pollDuration = pollDuration
	}
}

// WithID sets worker id
func WithID(id string) TickerOption {
	return func(wt *Ticker) {
		wt.id = id
	}
}
