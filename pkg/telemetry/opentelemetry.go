package telemetry

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	nrotel "github.com/newrelic/opentelemetry-exporter-go/newrelic"
// 	"go.opentelemetry.io/contrib/samplers/probability/consistent"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
// 	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
// 	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
// 	"go.opentelemetry.io/otel/sdk/resource"
// 	sdktrace "go.opentelemetry.io/otel/sdk/trace"
// 	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
// )

// func setupOpenTelemetry(ctx context.Context, cfg Config) error {
// 	otelResource := resource.NewWithAttributes(
// 		semconv.SchemaURL,
// 		semconv.ServiceNameKey.String(cfg.ServiceName),
// 		semconv.ServiceVersionKey.String("v0.5.1"), // TODO make this fetch from config
// 	)

// 	sdktraceOptions := []sdktrace.TracerProviderOption{
// 		sdktrace.WithResource(otelResource),
// 		sdktrace.WithSampler(consistent.ProbabilityBased(cfg.SamplingFraction)),
// 	}

// 	if cfg.EnableOtelAgent {
// 		otlptraceClient := otlptracehttp.NewClient(
// 			otlptracehttp.WithEndpoint(cfg.OpenTelAgentAddr),
// 		)
// 		spanExporter, err := otlptrace.New(ctx, otlptraceClient)
// 		if err != nil {
// 			return fmt.Errorf("creating OTLP trace exporter: %w", err)
// 		}

// 		go func() {
// 			<-ctx.Done()
// 			_ = spanExporter.Shutdown(context.Background())
// 		}()

// 		sdktraceOptions = append(sdktraceOptions, sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(spanExporter)))
// 	}

// 	if cfg.EnableDebugTrace {
// 		spanExporter, err := stdouttrace.New(
// 			stdouttrace.WithWriter(os.Stdout),
// 			stdouttrace.WithPrettyPrint(),
// 		)
// 		if err != nil {
// 			return fmt.Errorf("creating stdout trace exporter: %w", err)
// 		}

// 		sdktraceOptions = append(sdktraceOptions, sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(spanExporter)))
// 	}

// 	tracerProvider := sdktrace.NewTracerProvider(
// 		sdktraceOptions...,
// 	)

// 	otel.SetTracerProvider(tracerProvider)

// 	return nil
// }

// // trace.ApplyConfig(trace.Config{
// // 	DefaultSampler: trace.ProbabilitySampler(cfg.SamplingFraction),
// // })

// // if cfg.EnableMemory || cfg.EnableCPU {
// // 	opts := runmetrics.RunMetricOptions{
// // 		EnableCPU:    cfg.EnableCPU,
// // 		EnableMemory: cfg.EnableMemory,
// // 	}
// // 	if err := runmetrics.Enable(opts); err != nil {
// // 		return err
// // 	}
// // }

// // if err := setupViews(cfg.EnableHTTPServerMetrics); err != nil {
// // 	return err
// // }

// // if cfg.EnableNewrelic {
// // 	nrAppName := cfg.ServiceName
// // 	if cfg.NewRelicAppName != "" {
// // 		nrAppName = cfg.NewRelicAppName
// // 	}
// // 	exporter, err := nrcensus.NewExporter(nrAppName, cfg.NewRelicAPIKey)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	view.RegisterExporter(exporter)
// // 	trace.RegisterExporter(exporter)
// // }

// // if cfg.EnableOtelAgent {
// // 	ocExporter, err := ocagent.NewExporter(
// // 		ocagent.WithServiceName(cfg.ServiceName),
// // 		ocagent.WithInsecure(),
// // 		ocagent.WithAddress(cfg.OpenTelAgentAddr),
// // 	)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	go func() {
// // 		<-ctx.Done()
// // 		_ = ocExporter.Stop()
// // 	}()
// // 	trace.RegisterExporter(ocExporter)
// // 	view.RegisterExporter(ocExporter)
// // }

// // pe, err := prometheus.NewExporter(prometheus.Options{
// // 	Namespace: cfg.ServiceName,
// // })
// // if err != nil {
// // 	return err
// // }
// // view.RegisterExporter(pe)
// // mux.Handle("/metrics", pe)

// // zpages.Handle(mux, "/debug")
// // return nil
// // }
