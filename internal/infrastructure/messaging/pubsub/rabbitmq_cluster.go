// Package pubsub provides pub/sub messaging implementations
package pubsub

import (
	"context"
)

// RabbitMQCluster provides RabbitMQ cluster operations
type RabbitMQCluster interface {
	// Publish publishes a message to an exchange
	Publish(ctx context.Context, exchange string, routingKey string, message []byte) error

	// Consume consumes messages from a queue
	Consume(ctx context.Context, queue string, handler func([]byte) error) error

	// Close closes the cluster connection
	Close() error
}
