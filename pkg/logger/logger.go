package logger

import "context"

// contextKey is an unexported type for context keys to avoid collisions with other packages
type contextKey string

const traceIDKey contextKey = "traceID"

// TraceID extracts the trace ID from the context. Returns empty string if not found.
func TraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// WithTraceID adds a trace ID to the context.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// Logger is a simple logging interface. Each method accepts a message and optional key-value pairs
// for structured logging.
type Logger interface {
	// Debug logs DEBUG level messages. Arguments are key-value pairs. Example:
	//  logger.Debug("Starting process", "processID", string(123))
	Debug(msg string, keysAndValues ...string)

	// Info logs INFO level messages. Arguments are key-value pairs. Example:
	//  logger.Info("User logged in", "userID", string(123))
	Info(msg string, keysAndValues ...string)

	// Warn logs WARN level messages. Arguments are key-value pairs. Example:
	//  logger.Warn("Potential issue detected", "userID", string(123))
	Warn(msg string, keysAndValues ...string)

	// Error logs ERROR level messages. Arguments are key-value pairs. Example:
	//  logger.Error("Operation failed", "error", err.Error())
	Error(msg string, keysAndValues ...string)

	// DebugWithCtx logs DEBUG level messages with context. It's up to the implementation to decide
	// what to do with the context (e.g., extract trace IDs). Arguments are key-value pairs.
	DebugWithCtx(ctx context.Context, msg string, keysAndValues ...string)

	// InfoWithCtx logs INFO level messages with context. It's up to the implementation to decide
	// what to do with the context (e.g., extract trace IDs). Arguments are key-value pairs.
	InfoWithCtx(ctx context.Context, msg string, keysAndValues ...string)

	// WarnWithCtx logs WARN level messages with context. It's up to the implementation to decide
	// what to do with the context (e.g., extract trace IDs). Arguments are key-value pairs.
	WarnWithCtx(ctx context.Context, msg string, keysAndValues ...string)

	// ErrorWithCtx logs ERROR level messages with context. It's up to the implementation to decide
	// what to do with the context (e.g., extract trace IDs). Arguments are key-value pairs.
	ErrorWithCtx(ctx context.Context, msg string, keysAndValues ...string)
}
