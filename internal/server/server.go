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
	"go.uber.org/zap"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	sirenv1beta1 "github.com/odpf/siren/internal/server/proto/odpf/siren/v1beta1"
	"github.com/odpf/siren/internal/server/v1beta1"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"

	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/odpf/salt/log"
	"golang.org/x/net/http2/h2c"
)

//go:embed siren.swagger.yaml
var swaggerFile embed.FS

type Config struct {
	Host string `mapstructure:"host" default:"localhost"`
	Port int    `mapstructure:"port" default:"8080"`
}

// RunServer runs the application server
func RunServer(
	c Config, logger log.Logger,
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
	sirenv1beta1.RegisterSirenServiceServer(grpcServer, v1beta1Server)

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
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	grpcConn, err := grpc.DialContext(timeoutGrpcDialCtx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	logger.Info("server is running", "host", c.Host, "port", c.Port)
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
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
