// Package events supplies placeholders for the event-sourcing subsystem.
package events

import "context"

// Store is a minimal event store stub.
type Store struct{}

// NewStore creates a stub event store.
func NewStore() *Store {
	return &Store{}
}

// Append is a no-op append operation.
func (s *Store) Append(context.Context, interface{}) error { return nil }
