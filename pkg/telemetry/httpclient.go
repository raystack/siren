package telemetry

// import (
// 	"context"
// 	"fmt"

// 	"go.opencensus.io/trace"
// )

// type HTTPClientSpan struct {
// 	host string
// 	port string
// }

// func InitHTTPClientSpan(host string, port string) *HTTPClientSpan {
// 	return &HTTPClientSpan{
// 		host: host,
// 		port: port,
// 	}
// }

// func (c HTTPClientSpan) StartSpan(ctx context.Context, method string, route string, url string, spanAttributes map[string]string) (context.Context, *trace.Span) {
// 	// Refer https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/http.md
// 	ctx, span := trace.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", method, route), trace.WithSpanKind(trace.SpanKindClient))

// 	traceAttributes := []trace.Attribute{
// 		trace.StringAttribute("http.method", method),
// 		trace.StringAttribute("http.url", route),
// 		// trace.StringAttribute("http.resend_count", tryCount),
// 		trace.StringAttribute("net.peer.name", c.host),
// 		trace.StringAttribute("net.peer.port", c.port),
// 	}

// 	for k, v := range spanAttributes {
// 		traceAttributes = append(traceAttributes, trace.StringAttribute(k, v))
// 	}

// 	span.AddAttributes(
// 		traceAttributes...,
// 	)

// 	return ctx, span
// }
