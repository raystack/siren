package notification

import "time"

type HandlerOption func(*Handler)

func HandlerWithPollDuration(pollDuration time.Duration) HandlerOption {
	return func(h *Handler) {
		h.pollDuration = pollDuration
	}
}

func HandlerWithBatchSize(bs int) HandlerOption {
	return func(h *Handler) {
		h.batchSize = bs
	}
}
