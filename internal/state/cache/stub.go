// Package cache keeps only placeholder implementations at the moment.
package cache

import "context"

// StateCache is a trivial cache stub.
type StateCache struct{}

// NewStateCache creates a no-op cache.
func NewStateCache() *StateCache {
	return &StateCache{}
}

// Set stores a value (ignored in stub).
func (c *StateCache) Set(context.Context, string, interface{}) error { return nil }

// Get retrieves a value (always misses in stub).
func (c *StateCache) Get(context.Context, string) (interface{}, bool) {
	return nil, false
}
