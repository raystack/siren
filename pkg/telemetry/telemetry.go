package telemetry

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/odpf/salt/log"
)

type Config struct {
	// Debug sets the bind address for pprof & zpages server.
	Debug string `mapstructure:"debug_addr" yaml:"debug_addr" default:"localhost:8081"`

	// OpenCensus trace & metrics configurations.
	EnableCPU               bool    `mapstructure:"enable_cpu" yaml:"enable_cpu" default:"true"`
	EnableMemory            bool    `mapstructure:"enable_memory" yaml:"enable_memory" default:"true"`
	SamplingFraction        float64 `mapstructure:"sampling_fraction" yaml:"sampling_fraction" default:"1"`
	EnableHTTPServerMetrics bool    `mapstructure:"enable_http_server_metrics" yaml:"enable_http_server_metrics" default:"false"`

	// OpenCensus exporter configurations.
	ServiceName string `mapstructure:"service_name" yaml:"service_name" default:"siren"`

	// NewRelic exporter.
	EnableNewrelic  bool   `mapstructure:"enable_newrelic" yaml:"enable_newrelic" default:"false"`
	NewRelicAppName string `mapstructure:"newrelic_app_name"`
	NewRelicAPIKey  string `mapstructure:"newrelic_api_key" yaml:"newrelic_api_key" default:"____LICENSE_STRING_OF_40_CHARACTERS_____"`

	// OpenTelemetry Agent exporter.
	EnableOtelAgent  bool   `mapstructure:"enable_otel_agent" yaml:"enable_otel_agent" default:"false"`
	OpenTelAgentAddr string `mapstructure:"otel_agent_addr" yaml:"otel_agent_addr" default:"localhost:8088"`
	// EnableDebugTrace bool   `mapstructure:"enable_debug_trace"`
}

// Init initialises OpenCensus based async-telemetry processes and
// returns (i.e., it does not block).
func Init(ctx context.Context, cfg Config, lg log.Logger) {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))

	if err := setupOpenCensus(ctx, mux, cfg); err != nil {
		lg.Error("failed to setup OpenCensus", "err", err.Error())
	}

	if cfg.Debug != "" {
		go func() {
			if err := http.ListenAndServe(cfg.Debug, mux); err != nil {
				lg.Error("debug server exited due to error", "err", err.Error())
			}
		}()
	}
}
