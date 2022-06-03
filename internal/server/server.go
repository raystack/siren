package server

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime/middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/odpf/salt/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1"
	"github.com/odpf/siren/pkg/telemetry"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"

	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/internal/store"
	"golang.org/x/net/http2/h2c"
)

//go:embed siren.swagger.json
var swaggerFile embed.FS

// getZapLogLevelFromString helps to set logLevel from string
func getZapLogLevelFromString(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	nr, err := telemetry.New(&c.NewRelic)
	if err != nil {
		return err
	}

	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(c.Log.Level))
	logger := log.NewZap(log.ZapWithConfig(defaultConfig, zap.AddCaller()))
	zapper, err := zap.NewProduction()
	if err != nil {
		return err
	}

	gormDB, err := store.New(&c.DB)
	if err != nil {
		return err
	}

	httpClient := &http.Client{}
	repositories := store.NewRepositoryContainer(gormDB)
	services, err := v1beta1.InitContainer(repositories, gormDB, c, httpClient)
	if err != nil {
		return err
	}

	loggerOpts := []grpc_zap.Option{grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel)}

	// init grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nr),
			grpc_validator.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapper, loggerOpts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			nrgrpc.StreamServerInterceptor(nr),
			grpc_validator.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zapper),
		)),
	}
	grpcServer := grpc.NewServer(opts...)
	sirenv1beta1.RegisterSirenServiceServer(grpcServer, v1beta1.NewGRPCServer(services, nr, logger))

	// init http proxy
	timeoutGrpcDialCtx, grpcDialCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer grpcDialCancel()

	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	address := fmt.Sprintf(":%d", c.Port)
	grpcConn, err := grpc.DialContext(timeoutGrpcDialCtx, address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	runtimeCtx, runtimeCancel := context.WithCancel(context.Background())
	defer runtimeCancel()

	if err := sirenv1beta1.RegisterSirenServiceHandler(runtimeCtx, gwmux, grpcConn); err != nil {
		return err
	}

	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/siren.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(swaggerFile)).ServeHTTP(w, r)
	})
	baseMux.Handle("/documentation", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/siren.swagger.json",
		Path:    "documentation",
	}, http.NotFoundHandler()))
	baseMux.Handle("/", gwmux)

	server := &http.Server{
		Handler:      grpcHandlerFunc(grpcServer, baseMux),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("server is running", "port", c.Port)
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func RunMigrations(c *domain.Config) error {
	gormDB, err := store.New(&c.DB)
	if err != nil {
		return err
	}

	if err != nil {
		return nil
	}
	httpClient := &http.Client{}
	repositories := store.NewRepositoryContainer(gormDB)
	services, err := v1beta1.InitContainer(repositories, gormDB, c, httpClient)
	if err != nil {
		return err
	}

	err = services.MigrateAll(gormDB)
	if err != nil {
		return err
	}
	return nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
