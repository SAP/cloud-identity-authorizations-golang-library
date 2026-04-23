package logging

import "context"

type LogLevel int

type Logger interface {
	Debug(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
}

func Default() Logger {
	return NopLogger{}
}

type NopLogger struct{}

func (l NopLogger) Debug(ctx context.Context, msg string) {}
func (l NopLogger) Info(ctx context.Context, msg string)  {}
func (l NopLogger) Warn(ctx context.Context, msg string)  {}
func (l NopLogger) Error(ctx context.Context, msg string) {}
