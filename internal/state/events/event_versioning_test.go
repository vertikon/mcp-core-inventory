package events

import (
	"context"
	"testing"
	"time"
)

func TestNewEventVersioning(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, nil)

	if versioning == nil {
		t.Fatal("NewEventVersioning returned nil")
	}
}

func TestEventVersioningImpl_GetVersion_IncrementVersion(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, nil)
	ctx := context.Background()

	// Save initial event
	event1 := &Event{
		ID:            "event-1",
		Type:          EventTypeCreate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       1,
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}
	_ = store.SaveEvent(ctx, event1)

	// Get version
	versionInfo, err := versioning.GetVersion(ctx, "agg-1")
	if err != nil {
		t.Fatalf("GetVersion() error = %v", err)
	}

	if versionInfo == nil {
		t.Fatal("GetVersion() returned nil")
	}

	if versionInfo.CurrentVersion != 1 {
		t.Errorf("Expected version 1, got %d", versionInfo.CurrentVersion)
	}

	// Increment version
	event2 := &Event{
		ID:            "event-2",
		Type:          EventTypeUpdate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       0, // Will be auto-incremented
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}

	newVersion, err := versioning.IncrementVersion(ctx, "agg-1", event2)
	if err != nil {
		t.Fatalf("IncrementVersion() error = %v", err)
	}

	if newVersion != 2 {
		t.Errorf("Expected version 2, got %d", newVersion)
	}
}

func TestEventVersioningImpl_ValidateVersion(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, nil)
	ctx := context.Background()

	// Save event
	event := &Event{
		ID:            "event-1",
		Type:          EventTypeCreate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       1,
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}
	_ = store.SaveEvent(ctx, event)

	// Validate correct version (should succeed)
	err := versioning.ValidateVersion(ctx, "agg-1", 1)
	if err != nil {
		t.Fatalf("ValidateVersion() error = %v", err)
	}

	// Validate incorrect version (should fail)
	err = versioning.ValidateVersion(ctx, "agg-1", 2)
	if err == nil {
		t.Error("ValidateVersion() should fail for incorrect version")
	}
}

func TestEventVersioningImpl_GetVersionHistory(t *testing.T) {
	config := DefaultVersioningConfig()
	config.EnableHistory = true
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, config)
	ctx := context.Background()

	// Save and increment versions
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	for _, event := range events {
		_, _ = versioning.IncrementVersion(ctx, "agg-1", event)
	}

	// Get version history
	history, err := versioning.GetVersionHistory(ctx, "agg-1", 0)
	if err != nil {
		t.Fatalf("GetVersionHistory() error = %v", err)
	}

	if len(history) < 2 {
		t.Errorf("Expected >= 2 history entries, got %d", len(history))
	}
}

func TestEventVersioningImpl_ResolveVersionConflict(t *testing.T) {
	config := DefaultVersioningConfig()
	config.ConflictResolution = "accept-higher"
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, config)
	ctx := context.Background()

	conflict := &VersionConflict{
		AggregateID:     "agg-1",
		ExpectedVersion: 1,
		ActualVersion:   2,
		ConflictTime:    time.Now(),
	}

	resolved, err := versioning.ResolveVersionConflict(ctx, conflict)
	if err != nil {
		t.Fatalf("ResolveVersionConflict() error = %v", err)
	}

	if resolved != 2 {
		t.Errorf("Expected resolved version 2, got %d", resolved)
	}
}

func TestEventVersioningImpl_GetVersioningStats(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	versioning := NewEventVersioning(store, nil)
	ctx := context.Background()

	// Increment versions
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	for _, event := range events {
		_, _ = versioning.IncrementVersion(ctx, "agg-1", event)
	}

	// Get stats
	stats, err := versioning.GetVersioningStats(ctx)
	if err != nil {
		t.Fatalf("GetVersioningStats() error = %v", err)
	}

	if stats == nil {
		t.Fatal("GetVersioningStats() returned nil")
	}

	if stats.TotalVersions == 0 {
		t.Error("TotalVersions should be > 0")
	}
}

func TestDefaultVersioningConfig(t *testing.T) {
	config := DefaultVersioningConfig()

	if config == nil {
		t.Fatal("DefaultVersioningConfig returned nil")
	}

	if config.Strategy == "" {
		t.Error("Strategy should not be empty")
	}
}
