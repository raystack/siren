package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/goto/salt/log"

	"github.com/goto/siren/pkg/errors"
)

const (
	defaultBatchSize = 1
)

// Handler is a process to handle message publishing
type Handler struct {
	logger                 log.Logger
	q                      Queuer
	identifier             string
	notifierRegistry       map[string]Notifier
	supportedReceiverTypes []string
	batchSize              int
}

// NewHandler creates a new handler with some supported type of Notifiers
func NewHandler(cfg HandlerConfig, logger log.Logger, q Queuer, registry map[string]Notifier, opts ...HandlerOption) *Handler {
	h := &Handler{
		batchSize: defaultBatchSize,

		logger:           logger,
		notifierRegistry: registry,
		q:                q,
	}

	if cfg.BatchSize != 0 {
		h.batchSize = cfg.BatchSize
	}
	registeredReceivers := make([]string, 0, len(h.notifierRegistry))
	for k := range h.notifierRegistry {
		registeredReceivers = append(registeredReceivers, k)
	}
	h.supportedReceiverTypes = registeredReceivers

	if len(cfg.ReceiverTypes) != 0 {
		newSupportedReceiverTypes := []string{}
		for _, rt := range cfg.ReceiverTypes {
			found := false
			for _, k := range registeredReceivers {
				if rt == k {
					found = true
					break
				}
			}
			if found {
				newSupportedReceiverTypes = append(newSupportedReceiverTypes, rt)
			}
		}
		h.supportedReceiverTypes = newSupportedReceiverTypes
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handler) getNotifierPlugin(receiverType string) (Notifier, error) {
	receiverPlugin, exist := h.notifierRegistry[receiverType]
	if !exist {
		return nil, errors.ErrInvalid.WithMsgf("unsupported receiver type: %q on handler %s", receiverType, h.identifier)
	}
	return receiverPlugin, nil
}

func (h *Handler) Process(ctx context.Context, runAt time.Time) error {
	receiverTypes := h.supportedReceiverTypes
	if len(receiverTypes) == 0 {
		return errors.New("no receiver type plugin registered, skipping dequeue")
	} else {
		if err := h.q.Dequeue(ctx, receiverTypes, h.batchSize, h.MessageHandler); err != nil {
			if !errors.Is(err, ErrNoMessage) {
				return fmt.Errorf("dequeue failed on handler with id %s: %w", h.identifier, err)
			}
		}
	}
	return nil
}

// MessageHandler is a function to handler dequeued message
func (h *Handler) MessageHandler(ctx context.Context, messages []Message) error {
	for _, message := range messages {
		notifier, err := h.getNotifierPlugin(message.ReceiverType)
		if err != nil {
			return err
		}

		message.MarkPending(time.Now())

		newConfig, err := notifier.PostHookQueueTransformConfigs(ctx, message.Configs)
		if err != nil {
			message.MarkFailed(time.Now(), false, err)

			if cerr := h.q.ErrorCallback(ctx, message); cerr != nil {
				return cerr
			}
			return err
		}
		message.Configs = newConfig

		if retryable, err := notifier.Send(ctx, message); err != nil {
			message.MarkFailed(time.Now(), retryable, err)

			if cerr := h.q.ErrorCallback(ctx, message); cerr != nil {
				return cerr
			}
			return err
		}

		message.MarkPublished(time.Now())

		if err := h.q.SuccessCallback(ctx, message); err != nil {
			return err
		}
	}

	return nil
}
