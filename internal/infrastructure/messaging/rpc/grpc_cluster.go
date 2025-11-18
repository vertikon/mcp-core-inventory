// Package rpc provides RPC implementations
package rpc

import (
	"context"
)

// GRPCCluster provides gRPC cluster operations
type GRPCCluster interface {
	// Call makes a gRPC call
	Call(ctx context.Context, service string, method string, request interface{}) (interface{}, error)

	// Close closes the cluster connection
	Close() error
}
