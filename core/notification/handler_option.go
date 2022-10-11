package notification

import "time"

// HandlerOption is an option to customize handler creation
type HandlerOption func(*Handler)

// HandlerWithPollDuration sets created handler with the specified poll duration
func HandlerWithPollDuration(pollDuration time.Duration) HandlerOption {
	return func(h *Handler) {
		h.pollDuration = pollDuration
	}
}

// HandlerWithBatchSize sets created handler with the specified batch size
func HandlerWithBatchSize(bs int) HandlerOption {
	return func(h *Handler) {
		h.batchSize = bs
	}
}
