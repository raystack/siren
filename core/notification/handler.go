package notification

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/pkg/errors"
)

const (
	defaultPollDuration = 5 * time.Second
	defaultBatchSize    = 1
)

// Handler is a process to handle message publishing
type Handler struct {
	id                     string
	logger                 log.Logger
	q                      Queuer
	notifierRegistry       map[string]Notifier
	supportedReceiverTypes []string

	batchSize    int
	pollDuration time.Duration
}

// NewHandler creates a new handler with some supported type of Notifiers
func NewHandler(logger log.Logger, q Queuer, registry map[string]Notifier, opts ...HandlerOption) *Handler {
	h := &Handler{
		id:           uuid.NewString(),
		batchSize:    defaultBatchSize,
		pollDuration: defaultPollDuration,

		logger:           logger,
		notifierRegistry: registry,
		q:                q,
	}

	keys := make([]string, 0, len(h.notifierRegistry))
	for k := range h.notifierRegistry {
		keys = append(keys, k)
	}
	h.supportedReceiverTypes = keys

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) getNotifierPlugin(receiverType string) (Notifier, error) {
	receiverPlugin, exist := h.notifierRegistry[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q", receiverType)
	}
	return receiverPlugin, nil
}

func (h *Handler) RunHandler(ctx context.Context, wg *sync.WaitGroup, cancelChan chan struct{}) {
	defer wg.Done()

	ticker := time.NewTicker(h.pollDuration)
	defer ticker.Stop()

	h.logger.Info("running handler", "id", h.id)

	for {
		select {
		case <-cancelChan:
			h.logger.Info("stopping handler", "id", h.id)
			return

		case t := <-ticker.C:
			receiverTypes := h.supportedReceiverTypes
			if len(receiverTypes) == 0 {
				h.logger.Warn("no receiver type plugin registered, skipping dequeue", "scope", "notification.handler")
			} else {
				h.logger.Debug("dequeueing and publishing messages", "scope", "notification.handler", "receivers", receiverTypes, "batch size", h.batchSize, "running_at", t)
				if err := h.q.Dequeue(ctx, receiverTypes, h.batchSize, h.MessageHandler); err != nil && err != ErrNoMessage {
					h.logger.Error("dequeue failed", "scope", "notification.handler", "error", err)
				}
			}
		}
	}
}

// MessageHandler is a function to handler dequeued message
func (h *Handler) MessageHandler(ctx context.Context, messages []Message) error {
	for _, message := range messages {
		notifier, err := h.getNotifierPlugin(message.ReceiverType)
		if err != nil {
			return err
		}

		message.MarkPending(time.Now())

		if retryable, err := notifier.Publish(ctx, message); err != nil {

			message.MarkFailed(time.Now(), retryable, err)

			if err := h.q.ErrorHandler(ctx, message); err != nil {
				return err
			}
			return err
		}

		message.MarkPublished(time.Now())

		if err := h.q.SuccessHandler(ctx, message); err != nil {
			return err
		}

	}

	return nil
}
