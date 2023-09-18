package zaputil

import (
	"github.com/goto/salt/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(serviceName string, logLevel string, isGCPCompatible bool) log.Logger {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(logLevel))

	if isGCPCompatible {
		defaultConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(getZapLogLevelFromString(logLevel)),
			Encoding:    "json",
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			EncoderConfig:    EncoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	return log.NewZap(log.ZapWithConfig(
		defaultConfig,
		zap.Fields(ServiceContext(serviceName)),
		zap.AddCaller(),
		zap.AddStacktrace(zap.DPanicLevel),
	))
}

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
