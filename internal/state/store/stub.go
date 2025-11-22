// Package store exposes a stub implementation of the distributed store.
package store

import "context"

// DistributedStore represents a stubbed store.
type DistributedStore struct{}

// NewDistributedStore creates a new stub store.
func NewDistributedStore() *DistributedStore {
	return &DistributedStore{}
}

// Replicate performs a no-op replication.
func (s *DistributedStore) Replicate(context.Context) error { return nil }
