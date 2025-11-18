// Package logger provides structured logging using zap with trace_id and span_id support.
package logger

// LogLevel represents a log level
type LogLevel string

const (
	// LevelDebug is the debug log level
	LevelDebug LogLevel = "debug"
	// LevelInfo is the info log level
	LevelInfo LogLevel = "info"
	// LevelWarn is the warn log level
	LevelWarn LogLevel = "warn"
	// LevelError is the error log level
	LevelError LogLevel = "error"
)

// SetLevel sets the global log level
func SetLevel(level LogLevel) error {
	return Init(string(level), false)
}
