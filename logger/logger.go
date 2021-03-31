package logger

import (
	"github.com/odpf/siren/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(config *domain.LogConfig) (*zap.Logger, error) {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(config.Level))
	logger, err := zap.NewProductionConfig().Build()
	return logger, err
}

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
