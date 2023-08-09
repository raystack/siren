package telemetry

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opencensus.io/trace"
)

type MessagingTracer struct {
	queueSystem       string
	nrProducerSegment newrelic.MessageProducerSegment
	span              *trace.Span
}

func NewMessagingTracer(queueSystem string) *MessagingTracer {
	return &MessagingTracer{
		queueSystem: queueSystem,
	}
}

func (msg *MessagingTracer) StartSpan(ctx context.Context, op string, spanAttributes ...trace.Attribute) (context.Context, *trace.Span) {
	nrTx := newrelic.FromContext(ctx)
	msg.nrProducerSegment = newrelic.MessageProducerSegment{
		Library:              msg.queueSystem,
		DestinationType:      newrelic.MessageExchange,
		DestinationName:      "notification_queue",
		DestinationTemporary: false,
		StartTime:            nrTx.StartSegmentNow(),
	}

	// Refer https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md
	ctx, span := trace.StartSpan(ctx, fmt.Sprintf("notification_queue %s", op), trace.WithSpanKind(trace.SpanKindClient))

	traceAttributes := []trace.Attribute{
		trace.StringAttribute("messaging.system", msg.queueSystem),
		trace.StringAttribute("messaging.destination", "notification_queue"),
		trace.StringAttribute("messaging.destination_kind", "queue"),
		trace.StringAttribute("messaging.operation", op),
	}

	traceAttributes = append(traceAttributes, spanAttributes...)

	span.AddAttributes(
		traceAttributes...,
	)

	msg.span = span

	return ctx, span
}

func (msg *MessagingTracer) StopSpan() {
	msg.nrProducerSegment.End()
	msg.span.End()
}
