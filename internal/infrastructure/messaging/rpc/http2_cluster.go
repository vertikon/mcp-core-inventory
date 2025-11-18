// Package rpc provides RPC implementations
package rpc

import (
	"context"
)

// HTTP2Cluster provides HTTP/2 cluster operations
type HTTP2Cluster interface {
	// Call makes an HTTP/2 call
	Call(ctx context.Context, endpoint string, request interface{}) (interface{}, error)

	// Close closes the cluster connection
	Close() error
}
