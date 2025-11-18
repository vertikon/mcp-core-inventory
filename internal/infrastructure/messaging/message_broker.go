// Package messaging provides message broker abstractions
package messaging

import (
	"context"
)

// MessageBroker provides unified interface for streaming/pubsub
type MessageBroker interface {
	// Publish publishes a message to a topic
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe subscribes to a topic
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error

	// Close closes the broker connection
	Close() error
}
