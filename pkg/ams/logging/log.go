package logging

import "context"

type LogLevel int

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

func Default() Logger {
	return NoopLogger{}
}

type NoopLogger struct{}

func (l NoopLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (l NoopLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (l NoopLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (l NoopLogger) Errorf(ctx context.Context, format string, args ...interface{}) {}
