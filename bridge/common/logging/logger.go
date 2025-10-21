package logging

import (
	"context"
	"errors"
	"fmt"
)

// Log Level
const (
	LogLevelTrace uint8 = 1 << iota // 1   0000001
	LogLevelDebug                   // 2   0000010
	LogLevelInfo                    // 4   0000100
	LogLevelWarn                    // 8   0001000
	LogLevelError                   // 16  0010000
	LogLevelFatal                   // 32  0100000
	LogLevelPanic                   // 64  1000000

)

func LevelToString(level uint8) string {
	switch level {
	case LogLevelTrace:
		return "trace"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warning"
	case LogLevelError:
		return "error"
	case LogLevelFatal:
		return "fatal"
	case LogLevelPanic:
		return "panic"
	default:
		return "unknown"
	}
}

func LevelToNumber(level string) (uint8, error) {
	switch level {
	case "trace":
		return LogLevelTrace, nil
	case "debug":
		return LogLevelDebug, nil
	case "info":
		return LogLevelInfo, nil
	case "warning":
		return LogLevelWarn, nil
	case "error":
		return LogLevelError, nil
	case "fatal":
		return LogLevelFatal, nil
	case "panic":
		return LogLevelPanic, nil
	default:
		errMsg := fmt.Sprintf("unknown log level: %s", level)
		return LogLevelTrace, errors.New(errMsg)
	}
}

type Logger interface {
	Trace(args ...any)
	Debug(args ...any)
	Info(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Panic(args ...any)
}

type ContextLogger interface {
	Logger
	TraceContext(ctx context.Context, args ...any)
	DebugContext(ctx context.Context, args ...any)
	InfoContext(ctx context.Context, args ...any)
	WarnContext(ctx context.Context, args ...any)
	ErrorContext(ctx context.Context, args ...any)
	FatalContext(ctx context.Context, args ...any)
	PanicContext(ctx context.Context, args ...any)
}

func NOP() ContextLogger {
	return (*nopLogger)(nil)
}

type nopLogger struct {
}

func (n *nopLogger) Trace(args ...any) {
}

func (n *nopLogger) Debug(args ...any) {
}

func (n *nopLogger) Info(args ...any) {
}

func (n *nopLogger) Warn(args ...any) {
}

func (n *nopLogger) Error(args ...any) {
}

func (n *nopLogger) Fatal(args ...any) {
}

func (n *nopLogger) Panic(args ...any) {
}

func (n *nopLogger) TraceContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) DebugContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) InfoContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) WarnContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) ErrorContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) FatalContext(ctx context.Context, args ...any) {
}

func (n *nopLogger) PanicContext(ctx context.Context, args ...any) {
}
