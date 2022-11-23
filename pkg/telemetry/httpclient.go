package telemetry

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type Transport struct {
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := t.base()
	span := trace.FromContext(req.Context())

	span.AddAttributes([]trace.Attribute{
		trace.StringAttribute("span.kind", "client"),
	}...)

	ctx := trace.NewContext(req.Context(), span)

	return rt.RoundTrip(req.WithContext(ctx))
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return &ochttp.Transport{
			Base:           t.Base,
			NewClientTrace: ochttp.NewSpanAnnotatingClientTrace,
		}
	}
	return &ochttp.Transport{
		Base:           http.DefaultTransport,
		NewClientTrace: ochttp.NewSpanAnnotatingClientTrace,
	}
}
