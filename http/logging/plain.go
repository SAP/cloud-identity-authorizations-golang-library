package logging

import (
	"context"
	"fmt"
	"time"
)

type PlainLogger struct{}

func (l PlainLogger) Debug(ctx context.Context, msg string) {
	fmt.Printf("[%s] %s: %s\n", "DEBUG", time.Now().Format(time.RFC3339), msg)
}
func (l PlainLogger) Info(ctx context.Context, msg string) {
	fmt.Printf("[%s] %s: %s\n", "INFO", time.Now().Format(time.RFC3339), msg)
}
func (l PlainLogger) Warn(ctx context.Context, msg string) {
	fmt.Printf("[%s] %s: %s\n", "WARN", time.Now().Format(time.RFC3339), msg)
}
func (l PlainLogger) Error(ctx context.Context, msg string) {
	fmt.Printf("[%s] %s: %s\n", "ERROR", time.Now().Format(time.RFC3339), msg)
}
