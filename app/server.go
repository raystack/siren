package app

import (
	"context"
	"embed"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"strings"
	"time"

	cortexClient "github.com/grafana/cortex-tools/pkg/client"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/odpf/siren/api/handlers/v1"
	pb "github.com/odpf/siren/api/proto/odpf/siren"
	"github.com/odpf/siren/logger"
	"github.com/odpf/siren/metric"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/service"
	"github.com/odpf/siren/store"
	"golang.org/x/net/http2/h2c"
)

//go:embed siren.swagger.json
var swaggerFile embed.FS

// RunServer runs the application server
func RunServer(c *domain.Config) error {
	nr, err := metric.New(&c.NewRelic)
	if err != nil {
		return err
	}

	logger, err := logger.New(&c.Log)
	if err != nil {
		return err
	}

	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}
	cortexConfig := cortexClient.Config{
		Address:         c.Cortex.Address,
		UseLegacyRoutes: false,
	}
	client, err := cortexClient.New(cortexConfig)
	if err != nil {
		return nil
	}
	httpClient := &http.Client{}
	services, err := service.Init(store, c, client, httpClient)
	if err != nil {
		return err
	}

	loggerOpts := []grpc_zap.Option{grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel)}

	// init grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, loggerOpts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
		)),
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterSirenServiceServer(grpcServer, v1.NewGRPCServer(services, nr, logger))

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

	if err := pb.RegisterSirenServiceHandler(runtimeCtx, gwmux, grpcConn); err != nil {
		return err
	}

	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/siren.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(swaggerFile)).ServeHTTP(w, r)
		return
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

	log.Println("server running on port:", c.Port)
	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func RunMigrations(c *domain.Config) error {
	store, err := store.New(&c.DB)
	if err != nil {
		return err
	}

	cortexConfig := cortexClient.Config{
		Address:         c.Cortex.Address,
		UseLegacyRoutes: false,
	}
	client, err := cortexClient.New(cortexConfig)
	if err != nil {
		return nil
	}
	httpClient := &http.Client{}
	services, err := service.Init(store, c, client, httpClient)
	if err != nil {
		return err
	}

	services.MigrateAll(store)
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
