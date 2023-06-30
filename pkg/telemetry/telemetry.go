package telemetry

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/raystack/salt/log"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

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

func IncrementInt64Counter(ctx context.Context, si64 *stats.Int64Measure, tagMutator ...tag.Mutator) {
	counterCtx := ctx

	if tagMutator != nil {
		counterCtx, _ = tag.New(ctx, tagMutator...)
	}

	stats.Record(counterCtx, si64.M(1))
}

func GaugeMillisecond(ctx context.Context, si64 *stats.Int64Measure, value int64, tagMutator ...tag.Mutator) {
	counterCtx := ctx

	if tagMutator != nil {
		counterCtx, _ = tag.New(ctx, tagMutator...)
	}

	stats.Record(counterCtx, si64.M(value))
}
