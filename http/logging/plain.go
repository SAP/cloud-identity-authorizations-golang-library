package logging

import (
	"context"
	"fmt"
	"time"
)

type PlainLogger struct{}

func (l PlainLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	fmt.Printf("[%s] %s: %s\n", "DEBUG", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
func (l PlainLogger) Infof(ctx context.Context, msg string, args ...interface{}) {
	fmt.Printf("[%s] %s: %s\n", "INFO", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
func (l PlainLogger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	fmt.Printf("[%s] %s: %s\n", "WARN", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
func (l PlainLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	fmt.Printf("[%s] %s: %s\n", "ERROR", time.Now().Format(time.RFC3339), fmt.Sprintf(msg, args...))
}
