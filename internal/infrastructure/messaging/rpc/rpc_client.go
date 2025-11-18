// Package rpc provides RPC implementations
package rpc

import (
	"context"
)

// RPCClient provides generic RPC interface
type RPCClient interface {
	// Call makes an RPC call
	Call(ctx context.Context, service string, method string, request interface{}) (interface{}, error)

	// Close closes the client connection
	Close() error
}
