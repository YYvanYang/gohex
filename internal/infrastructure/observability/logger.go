package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger(env string) (*Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapLogger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		zap: zapLogger,
	}, nil
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.zap.Sugar().Infow(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.zap.Sugar().Errorw(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.zap.Sugar().Debugw(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.zap.Sugar().Warnw(msg, fields...)
} 