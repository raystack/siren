package telemetry

import (
	"context"
	"net/http"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/newrelic/newrelic-opencensus-exporter-go/nrcensus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/runmetrics"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

func setupOpenCensus(ctx context.Context, mux *http.ServeMux, cfg Config) error {
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(cfg.SamplingFraction),
	})

	if cfg.EnableMemory || cfg.EnableCPU {
		opts := runmetrics.RunMetricOptions{
			EnableCPU:    cfg.EnableCPU,
			EnableMemory: cfg.EnableMemory,
		}
		if err := runmetrics.Enable(opts); err != nil {
			return err
		}
	}

	if err := setupViews(); err != nil {
		return err
	}

	if cfg.EnableNewrelic {
		exporter, err := nrcensus.NewExporter(cfg.ServiceName, cfg.NewRelicAPIKey)
		if err != nil {
			return err
		}
		view.RegisterExporter(exporter)
		trace.RegisterExporter(exporter)
	}

	if cfg.EnableOtelAgent {
		ocExporter, err := ocagent.NewExporter(
			ocagent.WithServiceName(cfg.ServiceName),
			ocagent.WithInsecure(),
			ocagent.WithAddress(cfg.OpenTelAgentAddr),
		)
		if err != nil {
			return err
		}
		go func() {
			<-ctx.Done()
			_ = ocExporter.Stop()
		}()
		trace.RegisterExporter(ocExporter)
		view.RegisterExporter(ocExporter)
	}

	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: cfg.ServiceName,
	})
	if err != nil {
		return err
	}
	view.RegisterExporter(pe)
	mux.Handle("/metrics", pe)

	zpages.Handle(mux, "/debug")
	return nil
}

func setupViews() error {
	if err := setupDBViews(); err != nil {
		return err
	}

	if err := setupHTTPClientViews(); err != nil {
		return err
	}

	if err := setupGRPCServerViews(); err != nil {
		return err
	}

	if err := setupApplicationViews(); err != nil {
		return err
	}

	return nil
}

func setupGRPCServerViews() error {
	return view.Register(
		ocgrpc.ServerSentBytesPerRPCView,
		ocgrpc.ServerReceivedBytesPerRPCView,
		ocgrpc.ServerSentMessagesPerRPCView,
		ocgrpc.ServerReceivedMessagesPerRPCView,
		ocgrpc.ServerStartedRPCsView,
		ocgrpc.ServerCompletedRPCsView,
		ocgrpc.ServerLatencyView,
	)
}

func setupHTTPClientViews() error {
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		return err
	}

	clientHTTPViewTags := []tag.Key{
		ochttp.KeyClientMethod,
		ochttp.KeyClientStatus,
		ochttp.KeyClientHost,
	}

	return view.Register(
		&view.View{
			Name:        "opencensus.io/http/client/roundtrip_latency",
			Measure:     ochttp.ClientRoundtripLatency,
			Aggregation: ochttp.DefaultLatencyDistribution,
			Description: "End-to-end latency, by HTTP method and response status",
			TagKeys:     clientHTTPViewTags,
		},
		&view.View{
			Name:        "opencensus.io/http/client/sent_bytes",
			Measure:     ochttp.ClientSentBytes,
			Aggregation: ochttp.DefaultSizeDistribution,
			Description: "Total bytes sent in request body (not including headers), by HTTP method and response status",
			TagKeys:     clientHTTPViewTags,
		},
		&view.View{
			Name:        "opencensus.io/http/client/received_bytes",
			Measure:     ochttp.ClientReceivedBytes,
			Aggregation: ochttp.DefaultSizeDistribution,
			Description: "Total bytes received in response bodies (not including headers but including error responses with bodies), by HTTP method and response status",
			TagKeys:     clientHTTPViewTags,
		},
	)
}

func setupDBViews() error {
	return view.Register(ocsql.DefaultViews...)
}
