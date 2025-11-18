// Package streaming provides streaming message implementations
package streaming

import (
	"context"
)

// StreamClient provides generic streaming interface
type StreamClient interface {
	// Produce produces a message to a stream
	Produce(ctx context.Context, stream string, message []byte) error

	// Consume consumes messages from a stream
	Consume(ctx context.Context, stream string, handler func([]byte) error) error

	// Close closes the client connection
	Close() error
}
