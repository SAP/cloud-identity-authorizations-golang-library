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
	return NopLogger{}
}

type NopLogger struct{}

func (l NopLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (l NopLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (l NopLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (l NopLogger) Errorf(ctx context.Context, format string, args ...interface{}) {}
