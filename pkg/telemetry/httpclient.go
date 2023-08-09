package telemetry

import (
	"net/http"

	"go.opencensus.io/trace"
)

type Transport struct {
	Origin http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	span := trace.FromContext(req.Context())

	span.AddAttributes([]trace.Attribute{
		trace.StringAttribute("span.kind", "client"),
	}...)

	ctx := trace.NewContext(req.Context(), span)

	return t.Origin.RoundTrip(req.WithContext(ctx))
}
