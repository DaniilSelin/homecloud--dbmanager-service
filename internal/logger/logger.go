package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestID = "RequestID"
	LoggerKey = "logger"
)

type Logger struct {
	l *zap.Logger
}

type LoggerConfig struct {
	Level string
}

func New(level string) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	if level == "debug" {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	return &Logger{l: logger}, nil
}

func CtxWWithLogger(ctx context.Context, lg *Logger) context.Context {
	ctx = context.WithValue(ctx, LoggerKey, lg)
	return ctx
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(LoggerKey).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Info(msg, fields...)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Debug(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Error(msg, fields...)
}
