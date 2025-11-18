// Package pubsub provides pub/sub messaging implementations
package pubsub

import (
	"context"
)

// PubSubClient provides generic pub/sub interface
type PubSubClient interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error

	// Close closes the client connection
	Close() error
}
