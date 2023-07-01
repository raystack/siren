package zaputil

import (
	"github.com/raystack/salt/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

// GRPCZapLogger returns *zap.Logger from salt/log.Logger.
// if error, returns a new *zap.Logger instance
func GRPCZapLogger(logger log.Logger) (*zap.Logger, error) {
	var zapLogger *zap.Logger
	var err error
	zapLogger, err = zap.NewProduction()
	if err != nil {
		return nil, err
	}
	if mainZapLogger, ok := logger.(*log.Zap); !ok {
		logger.Warn("failed to get main logger to use in grpc interceptor, fallback to new logger")
	} else {
		zapLogger = mainZapLogger.GetInternalZapLogger().Desugar()
	}
	return zapLogger, nil
}

// GRPCCodeToLevel is the mapping of gRPC return codes and interceptor
// log level. Convert codes.OK to DEBUG level, the rest are
// the same with the DefaultCodeToLevel in [grpc_zap].
//
// [grpc_zap]: https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap#DefaultCodeToLevel
func GRPCCodeToLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.DebugLevel
	case codes.Canceled:
		return zap.InfoLevel
	case codes.Unknown:
		return zap.ErrorLevel
	case codes.InvalidArgument:
		return zap.InfoLevel
	case codes.DeadlineExceeded:
		return zap.WarnLevel
	case codes.NotFound:
		return zap.InfoLevel
	case codes.AlreadyExists:
		return zap.InfoLevel
	case codes.PermissionDenied:
		return zap.WarnLevel
	case codes.Unauthenticated:
		return zap.InfoLevel // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return zap.WarnLevel
	case codes.FailedPrecondition:
		return zap.WarnLevel
	case codes.Aborted:
		return zap.WarnLevel
	case codes.OutOfRange:
		return zap.WarnLevel
	case codes.Unimplemented:
		return zap.ErrorLevel
	case codes.Internal:
		return zap.ErrorLevel
	case codes.Unavailable:
		return zap.WarnLevel
	case codes.DataLoss:
		return zap.ErrorLevel
	default:
		return zap.ErrorLevel
	}
}
