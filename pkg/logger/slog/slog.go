package slog

import (
	"quote-service/pkg/logger"
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Logger struct {
	logger *slog.Logger
}

var _ logger.Logger = (*Logger)(nil)

type NewLoggerArgs struct {
	LogFormat string
}

func NewLogger(args NewLoggerArgs) *Logger {
	var logger *slog.Logger

	switch args.LogFormat {
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})) //nolint:exhaustruct
	case "color":
		logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{})) //nolint:exhaustruct
	}

	return &Logger{logger: logger}
}

// Debug implements logger.Logger.
func (l *Logger) Debug(msg string, keysAndValues ...string) {
	l.logger.Debug(msg, stringsToAnySlice(keysAndValues)...)
}

// DebugWithCtx implements logger.Logger.
func (l *Logger) DebugWithCtx(ctx context.Context, msg string, keysAndValues ...string) {
	l.logger.Debug(msg, l.withTraceID(ctx, keysAndValues)...)
}

// Error implements logger.Logger.
func (l *Logger) Error(msg string, keysAndValues ...string) {
	l.logger.Error(msg, stringsToAnySlice(keysAndValues)...)
}

// ErrorWithCtx implements logger.Logger.
func (l *Logger) ErrorWithCtx(ctx context.Context, msg string, keysAndValues ...string) {
	l.logger.Error(msg, l.withTraceID(ctx, keysAndValues)...)
}

// Info implements logger.Logger.
func (l *Logger) Info(msg string, keysAndValues ...string) {
	l.logger.Info(msg, stringsToAnySlice(keysAndValues)...)
}

// InfoWithCtx implements logger.Logger.
func (l *Logger) InfoWithCtx(ctx context.Context, msg string, keysAndValues ...string) {
	l.logger.Info(msg, l.withTraceID(ctx, keysAndValues)...)
}

// Warn implements logger.Logger.
func (l *Logger) Warn(msg string, keysAndValues ...string) {
	l.logger.Warn(msg, stringsToAnySlice(keysAndValues)...)
}

// WarnWithCtx implements logger.Logger.
func (l *Logger) WarnWithCtx(ctx context.Context, msg string, keysAndValues ...string) {
	l.logger.Warn(msg, l.withTraceID(ctx, keysAndValues)...)
}

// withTraceID extracts trace ID from context and prepends it to key-value pairs
func (l *Logger) withTraceID(ctx context.Context, keysAndValues []string) []any {
	traceID := logger.TraceID(ctx)
	if traceID == "" {
		return stringsToAnySlice(keysAndValues)
	}

	// Prepend traceID to the key-value pairs
	args := make([]any, 0, len(keysAndValues)+2) //nolint:mnd
	args = append(args, "traceID", traceID)
	args = append(args, stringsToAnySlice(keysAndValues)...)

	return args
}

func stringsToAnySlice(strs []string) []any {
	anys := make([]any, len(strs))
	for i, s := range strs {
		anys[i] = s
	}

	return anys
}
