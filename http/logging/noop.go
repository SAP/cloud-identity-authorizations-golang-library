package logging

import "context"

type NoopLogger struct{}

func (l NoopLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {}
func (l NoopLogger) Infof(ctx context.Context, msg string, args ...interface{})  {}
func (l NoopLogger) Warnf(ctx context.Context, msg string, args ...interface{})  {}
func (l NoopLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {}
