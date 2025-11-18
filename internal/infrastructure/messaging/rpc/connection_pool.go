// Package rpc provides RPC connection pool implementations
package rpc

import (
	"context"
)

// ConnectionPool provides connection pool operations
type ConnectionPool interface {
	// Get gets a connection from the pool
	Get(ctx context.Context) (interface{}, error)

	// Put returns a connection to the pool
	Put(conn interface{}) error

	// Close closes all connections in the pool
	Close() error
}
