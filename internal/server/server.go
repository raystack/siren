package server

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-openapi/runtime/middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/odpf/salt/log"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
)

//go:embed siren.swagger.yaml
var swaggerFile embed.FS

type Config struct {
	Host string `mapstructure:"host" default:"localhost"`
	Port int    `mapstructure:"port" default:"8080"`
}

func (cfg Config) addr() string { return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port) }

// RunServer runs the application server
func RunServer(
	c Config,
	logger log.Logger,
	nr *newrelic.Application,
	templateService v1beta1.TemplateService,
	ruleService v1beta1.RuleService,
	alertService v1beta1.AlertService,
	providerService v1beta1.ProviderService,
	namespaceService v1beta1.NamespaceService,
	receiverService v1beta1.ReceiverService,
	subscriptionService v1beta1.SubscriptionService) error {

	v1beta1Server := v1beta1.NewGRPCServer(
		nr,
		logger,
		templateService,
		ruleService,
		alertService,
		providerService,
		namespaceService,
		receiverService,
		subscriptionService,
	)

	// TODO grpc should uses the same log
	loggerOpts := []grpc_zap.Option{grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel)}
	zapper, err := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return err
	}

	// init grpc server
	grpcServer := grpc.NewServer(
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
	)

	sirenv1beta1.RegisterSirenServiceServer(grpcServer, v1beta1Server)
	grpc_health_v1.RegisterHealthServer(grpcServer, v1beta1Server)

	// init http proxy
	grpcDialCtx, grpcDialCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer grpcDialCancel()

	grpcConn, err := grpc.DialContext(grpcDialCtx, c.addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	runtimeCtx, runtimeCancel := context.WithCancel(context.Background())
	defer runtimeCancel()

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
		runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(grpcConn), "/ping"),
	)

	if err := sirenv1beta1.RegisterSirenServiceHandler(runtimeCtx, gwmux, grpcConn); err != nil {
		return err
	}

	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/siren.swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(swaggerFile)).ServeHTTP(w, r)
	})
	baseMux.Handle("/documentation", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/siren.swagger.yaml",
		Path:    "documentation",
	}, http.NotFoundHandler()))
	baseMux.Handle("/", gwmux)

	httpServer := &http.Server{
		Handler:      grpcHandlerFunc(grpcServer, baseMux),
		Addr:         c.addr(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("server is running", "host", c.Host, "port", c.Port)
	idleConnsClosed := make(chan struct{})
	interrupt := make(chan os.Signal, 1)
	go func() {
		signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt

		if httpServer != nil {
			// We received an interrupt signal, shut down.
			logger.Warn("stopping http server...")
			if err := httpServer.Shutdown(context.Background()); err != nil {
				logger.Error("HTTP server Shutdown", "err", err)
			}
		}

		if grpcServer != nil {
			logger.Warn("stopping grpc server...")
			grpcServer.GracefulStop()
		}

		// Close DB here if any

		close(idleConnsClosed)
	}()

	go func() {
		defer func() { interrupt <- syscall.SIGTERM }()
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("HTTP server ListenAndServe", "err", err)
		}
	}()

	logger.Info("server started")

	<-idleConnsClosed

	logger.Info("server stopped")

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
