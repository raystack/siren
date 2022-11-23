package telemetry

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
)

type MessagingSpan struct {
	queueSystem     string
	destinationName string
}

func InitMessagingSpan(queueSystem string, destinationName string) *MessagingSpan {
	return &MessagingSpan{
		queueSystem:     queueSystem,
		destinationName: destinationName,
	}
}

func (msg MessagingSpan) StartSpan(ctx context.Context, op string, messageID string, spanAttributes map[string]string) (context.Context, *trace.Span) {
	// Refer https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md
	ctx, span := trace.StartSpan(ctx, fmt.Sprintf("%s %s", msg.destinationName, op), trace.WithSpanKind(trace.SpanKindClient))

	traceAttributes := []trace.Attribute{
		trace.StringAttribute("messaging.system", msg.queueSystem),
		trace.StringAttribute("messaging.destination", msg.destinationName),
		trace.StringAttribute("messaging.destination_kind", "queue"),
		trace.StringAttribute("messaging.operation", op),
		trace.StringAttribute("messaging.message_id", messageID),
	}

	for k, v := range spanAttributes {
		traceAttributes = append(traceAttributes, trace.StringAttribute(k, v))
	}

	span.AddAttributes(
		traceAttributes...,
	)

	return ctx, span
}
