package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/goto/salt/log"
	"github.com/goto/salt/mux"
	"github.com/goto/siren/internal/api"
	"github.com/goto/siren/internal/api/v1beta1"
	"github.com/goto/siren/pkg/zaputil"
	swagger "github.com/goto/siren/proto"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const defaultGracePeriod = 5 * time.Second

type GRPCConfig struct {
	Port           int `yaml:"port" mapstructure:"port" default:"8081"`
	MaxRecvMsgSize int `yaml:"max_recv_msg_size" mapstructure:"max_recv_msg_size" default:"33554432"`
	MaxSendMsgSize int `yaml:"max_send_msg_size" mapstructure:"max_send_msg_size" default:"33554432"`
}

type Config struct {
	Host                  string            `mapstructure:"host" yaml:"host" default:"localhost"`
	Port                  int               `mapstructure:"port" yaml:"port" default:"8080"`
	EncryptionKey         string            `mapstructure:"encryption_key" yaml:"encryption_key" default:"_ENCRYPTIONKEY_OF_32_CHARACTERS_"`
	APIHeaders            api.HeadersConfig `mapstructure:"api_headers" yaml:"api_headers"`
	UseGlobalSubscription bool              `mapstructure:"use_global_subscription" yaml:"use_global_subscription" default:"false"`
	EnableSilenceFeature  bool              `mapstructure:"enable_silence_feature" yaml:"enable_silence_feature" default:"false"`
	DebugRequest          bool              `mapstructure:"debug_request" yaml:"debug_request" default:"false"`
	GRPC                  GRPCConfig        `mapstructure:"grpc"`
}

func (cfg Config) addr() string     { return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port) }
func (cfg Config) grpcAddr() string { return fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPC.Port) }

// RunServer runs the application server
func RunServer(
	ctx context.Context,
	c Config,
	logger log.Logger,
	apiDeps *api.Deps) error {

	var err error

	// init grpc server
	zapLogger, err := zaputil.GRPCZapLogger(logger)
	if err != nil {
		return err
	}

	loggerOpts := []grpc_zap.Option{
		grpc_zap.WithLevels(grpc_zap.DefaultCodeToLevel),
		grpc_zap.WithTimestampFormat(time.RFC3339Nano),
		grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && fullMethodName == grpc_health_v1.Health_Check_FullMethodName {
				return false
			}
			// by default everything will be logged
			return true
		}),
	}
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(zapLogger, loggerOpts...),
			grpc_recovery.UnaryServerInterceptor(),
			// interceptors.OtelsqlMeterInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zapLogger, loggerOpts...),
			grpc_recovery.StreamServerInterceptor(),
		)),
	)

	// init http proxy
	grpcDialCtx, grpcDialCancel := context.WithTimeout(ctx, time.Second*5)
	defer grpcDialCancel()

	grpcConn, err := grpc.DialContext(grpcDialCtx, c.grpcAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	httpGateway := runtime.NewServeMux(
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
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			return key, api.SupportedHeaders(c.APIHeaders)[key]
		}),
	)

	reflection.Register(grpcServer)

	runtimeCtx, runtimeCancel := context.WithCancel(ctx)
	defer runtimeCancel()

	sirenServiceRPC := v1beta1.NewGRPCServer(
		logger,
		c.APIHeaders,
		apiDeps,
		v1beta1.WithGlobalSubscription(c.UseGlobalSubscription),
		v1beta1.WithDebugRequest(c.DebugRequest),
	)
	grpcServer.RegisterService(&sirenv1beta1.SirenService_ServiceDesc, sirenServiceRPC)
	grpcServer.RegisterService(&grpc_health_v1.Health_ServiceDesc, sirenServiceRPC)
	if err := sirenv1beta1.RegisterSirenServiceHandler(runtimeCtx, httpGateway, grpcConn); err != nil {
		return err
	}

	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/siren.swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(swagger.File)).ServeHTTP(w, r)
	})
	baseMux.Handle("/documentation", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		SpecURL: "/siren.swagger.yaml",
		Path:    "documentation",
	}, http.NotFoundHandler()))
	baseMux.Handle("/", httpGateway)

	logger.Info("server is running", "host", c.Host, "port", c.Port, "grpc_port", c.GRPC.Port)

	return mux.Serve(runtimeCtx,
		mux.WithHTTPTarget(c.addr(), &http.Server{
			Handler:      baseMux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}),
		mux.WithGRPCTarget(c.grpcAddr(), grpcServer),
		mux.WithGracePeriod(defaultGracePeriod),
	)
}
