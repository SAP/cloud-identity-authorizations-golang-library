package log

type LogLevel int

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelError
)

type NopLogger struct{}

func (l NopLogger) Debugf(format string, args ...any) {}
func (l NopLogger) Infof(format string, args ...any)  {}
func (l NopLogger) Errorf(format string, args ...any) {}
