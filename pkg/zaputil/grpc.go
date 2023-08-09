package zaputil

import (
	"github.com/goto/salt/log"
	"go.uber.org/zap"
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
