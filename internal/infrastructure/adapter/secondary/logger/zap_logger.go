package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/gohex/gohex/internal/application/port"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(config LogConfig) (Logger, error) {
	var cfg zap.Config
	if config.Environment == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.OutputPaths = []string{config.OutputPath}
	cfg.Level = zap.NewAtomicLevelAt(getLogLevel(config.Level))

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: logger.Sugar(),
	}, nil
}

func (l *zapLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debugw(msg, args...)
}

func (l *zapLogger) Info(msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *zapLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warnw(msg, args...)
}

func (l *zapLogger) Error(msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *zapLogger) With(key string, value interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(key, value),
	}
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
} 