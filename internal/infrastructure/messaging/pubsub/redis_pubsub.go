// Package pubsub provides pub/sub messaging implementations
package pubsub

import (
	"context"
)

// RedisPubSub provides Redis pub/sub operations
type RedisPubSub interface {
	// Publish publishes a message to a channel
	Publish(ctx context.Context, channel string, message []byte) error

	// Subscribe subscribes to a channel
	Subscribe(ctx context.Context, channel string, handler func([]byte) error) error

	// Close closes the pub/sub connection
	Close() error
}
