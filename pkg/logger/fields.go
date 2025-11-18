// Package logger provides structured logging using zap with trace_id and span_id support.
package logger

import "go.uber.org/zap"

// Field helpers for common log fields

// String returns a zap.String field
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int returns a zap.Int field
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Error returns a zap.Error field
func ErrorField(err error) zap.Field {
	return zap.Error(err)
}

// Any returns a zap.Any field
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
