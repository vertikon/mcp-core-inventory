package pubsub

import "context"

// NATSPubSubClient is a placeholder implementation of PubSubClient backed by NATS.
type NATSPubSubClient struct{}

// NewNATSPubSubClient creates a new placeholder pub/sub client.
func NewNATSPubSubClient() PubSubClient {
	return &NATSPubSubClient{}
}

// Publish publishes a message. Real wiring should be done when infrastructure is available.
func (c *NATSPubSubClient) Publish(ctx context.Context, topic string, message []byte) error {
	return nil
}

// Subscribe registers a handler. The stub simply acknowledges the subscription.
func (c *NATSPubSubClient) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	return nil
}

// Close closes the client.
func (c *NATSPubSubClient) Close() error {
	return nil
}
