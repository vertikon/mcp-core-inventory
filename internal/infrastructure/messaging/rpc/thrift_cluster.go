// Package rpc provides RPC implementations
package rpc

import (
	"context"
)

// ThriftCluster provides Thrift cluster operations
type ThriftCluster interface {
	// Call makes a Thrift RPC call
	Call(ctx context.Context, service string, method string, request interface{}) (interface{}, error)

	// Close closes the cluster connection
	Close() error
}
