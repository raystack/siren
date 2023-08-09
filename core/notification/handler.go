package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/goto/salt/log"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"github.com/goto/siren/pkg/errors"
	"github.com/goto/siren/pkg/telemetry"
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
	messagingTracer        *telemetry.MessagingTracer
	nrApp                  *newrelic.Application

	batchSize int
}

// NewHandler creates a new handler with some supported type of Notifiers
func NewHandler(cfg HandlerConfig, logger log.Logger, nrApp *newrelic.Application, q Queuer, registry map[string]Notifier, opts ...HandlerOption) *Handler {
	h := &Handler{
		batchSize: defaultBatchSize,

		logger:           logger,
		nrApp:            nrApp,
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

	h.messagingTracer = telemetry.NewMessagingTracer(q.Type())

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
	txn := h.nrApp.StartTransaction(h.identifier)
	defer txn.End()
	nrCtx := newrelic.NewContext(ctx, txn)

	receiverTypes := h.supportedReceiverTypes
	if len(receiverTypes) == 0 {
		err := errors.New("no receiver type plugin registered, skipping dequeue")
		txn.NoticeError(err)
		return err
	} else {
		traceCtx, span := h.messagingTracer.StartSpan(nrCtx, "batch_dequeue", trace.StringAttribute("messaging.handler_id", h.identifier))
		defer h.messagingTracer.StopSpan()

		if err := h.q.Dequeue(traceCtx, receiverTypes, h.batchSize, h.MessageHandler); err != nil {
			if !errors.Is(err, ErrNoMessage) {
				span.SetStatus(trace.Status{
					Code:    trace.StatusCodeUnknown,
					Message: err.Error(),
				})
				err = fmt.Errorf("dequeue failed on handler with id %s: %w", h.identifier, err)
				txn.NoticeError(err)
				return err
			} else {
				// no messages found
				txn.Ignore()
			}
		}
	}
	return nil
}

// MessageHandler is a function to handler dequeued message
func (h *Handler) MessageHandler(ctx context.Context, messages []Message) error {
	for _, message := range messages {

		telemetry.GaugeMillisecond(ctx, telemetry.MetricNotificationMessageQueueTime, time.Since(message.UpdatedAt).Milliseconds(),
			tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

		notifier, err := h.getNotifierPlugin(message.ReceiverType)
		if err != nil {
			return err
		}

		message.MarkPending(time.Now())

		telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
			tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
			tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

		newConfig, err := notifier.PostHookQueueTransformConfigs(ctx, message.Configs)
		if err != nil {
			message.MarkFailed(time.Now(), false, err)

			telemetry.IncrementInt64Counter(ctx, telemetry.MetricReceiverHookFailed,
				tag.Upsert(telemetry.TagReceiverType, message.ReceiverType),
				tag.Upsert(telemetry.TagHookCondition, telemetry.HookConditionPostHookQueue),
			)

			telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
				tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
				tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

			if cerr := h.q.ErrorCallback(ctx, message); cerr != nil {
				return cerr
			}
			return err
		}
		message.Configs = newConfig

		if retryable, err := notifier.Send(ctx, message); err != nil {
			message.MarkFailed(time.Now(), retryable, err)

			telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
				tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
				tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

			if cerr := h.q.ErrorCallback(ctx, message); cerr != nil {
				return cerr
			}
			return err
		}

		message.MarkPublished(time.Now())

		telemetry.IncrementInt64Counter(ctx, telemetry.MetricNotificationMessageCounter,
			tag.Upsert(telemetry.TagMessageStatus, message.Status.String()),
			tag.Upsert(telemetry.TagReceiverType, message.ReceiverType))

		if err := h.q.SuccessCallback(ctx, message); err != nil {
			return err
		}
	}

	return nil
}
