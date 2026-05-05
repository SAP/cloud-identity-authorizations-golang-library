package logging

import "context"

type Logger interface {
	Debugf(ctx context.Context, msg string, args ...interface{})
	Infof(ctx context.Context, msg string, args ...interface{})
	Warnf(ctx context.Context, msg string, args ...interface{})
	Errorf(ctx context.Context, msg string, args ...interface{})
}
