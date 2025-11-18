// Package streaming provides streaming message implementations
package streaming

import (
	"context"
)

// PulsarCluster provides Pulsar cluster operations
type PulsarCluster interface {
	// Produce produces a message to a topic
	Produce(ctx context.Context, topic string, message []byte) error

	// Consume consumes messages from a topic
	Consume(ctx context.Context, topic string, handler func([]byte) error) error

	// Close closes the cluster connection
	Close() error
}
