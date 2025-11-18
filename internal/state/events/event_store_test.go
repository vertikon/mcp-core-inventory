package events

import (
	"context"
	"testing"
	"time"
)

func TestNewInMemoryEventStore(t *testing.T) {
	config := DefaultEventStoreConfig()
	store := NewInMemoryEventStore(config)

	if store == nil {
		t.Fatal("NewInMemoryEventStore returned nil")
	}

	if store.events == nil {
		t.Error("events map should not be nil")
	}
}

func TestInMemoryEventStore_SaveEvent_GetEvents(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	event := &Event{
		ID:            "event-1",
		Type:          EventTypeCreate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       1,
		Data:          map[string]interface{}{"key": "value"},
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}

	// Save event
	err := store.SaveEvent(ctx, event)
	if err != nil {
		t.Fatalf("SaveEvent() error = %v", err)
	}

	// Get events
	events, err := store.GetEvents(ctx, "agg-1", 1, 1)
	if err != nil {
		t.Fatalf("GetEvents() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].ID != "event-1" {
		t.Errorf("Expected event ID 'event-1', got '%s'", events[0].ID)
	}
}

func TestInMemoryEventStore_SaveEvents_Batch(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	events := []*Event{
		{
			ID:            "event-1",
			Type:          EventTypeCreate,
			AggregateID:   "agg-1",
			AggregateType: "test",
			Version:       1,
			Timestamp:     time.Now(),
			NodeID:        "node-1",
		},
		{
			ID:            "event-2",
			Type:          EventTypeUpdate,
			AggregateID:   "agg-1",
			AggregateType: "test",
			Version:       2,
			Timestamp:     time.Now(),
			NodeID:        "node-1",
		},
	}

	// Save events
	err := store.SaveEvents(ctx, events)
	if err != nil {
		t.Fatalf("SaveEvents() error = %v", err)
	}

	// Get all events
	allEvents, err := store.GetAllEvents(ctx, "agg-1")
	if err != nil {
		t.Fatalf("GetAllEvents() error = %v", err)
	}

	if len(allEvents) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(allEvents))
	}
}

func TestInMemoryEventStore_VersionContinuity(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	// Save first event
	event1 := &Event{
		ID:            "event-1",
		Type:          EventTypeCreate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       1,
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}
	err := store.SaveEvent(ctx, event1)
	if err != nil {
		t.Fatalf("SaveEvent() error = %v", err)
	}

	// Try to save event with wrong version (should fail)
	event2 := &Event{
		ID:            "event-2",
		Type:          EventTypeUpdate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       3, // Should be 2
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	}
	err = store.SaveEvent(ctx, event2)
	if err == nil {
		t.Error("SaveEvent() should fail for version gap")
	}
}

func TestInMemoryEventStore_StreamEvents(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	// Save some events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	// Stream events
	eventCh, err := store.StreamEvents(ctx, "agg-1", 1)
	if err != nil {
		t.Fatalf("StreamEvents() error = %v", err)
	}

	count := 0
	for event := range eventCh {
		if event == nil {
			break
		}
		count++
		if count >= 2 {
			break
		}
	}

	if count < 2 {
		t.Errorf("Expected at least 2 events, got %d", count)
	}
}

func TestInMemoryEventStore_CreateSnapshot_GetSnapshot(t *testing.T) {
	store := NewInMemoryEventStore(nil)
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

	// Create snapshot
	err := store.CreateSnapshot(ctx, "agg-1", 1, map[string]interface{}{"state": "snapshot"})
	if err != nil {
		t.Fatalf("CreateSnapshot() error = %v", err)
	}

	// Get snapshot
	snapshot, err := store.GetSnapshot(ctx, "agg-1")
	if err != nil {
		t.Fatalf("GetSnapshot() error = %v", err)
	}

	if snapshot == nil {
		t.Fatal("GetSnapshot() returned nil")
	}

	if snapshot.Version != 1 {
		t.Errorf("Expected version 1, got %d", snapshot.Version)
	}
}

func TestInMemoryEventStore_GetEventStats(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	// Save some events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	// Get stats
	stats, err := store.GetEventStats(ctx)
	if err != nil {
		t.Fatalf("GetEventStats() error = %v", err)
	}

	if stats == nil {
		t.Fatal("GetEventStats() returned nil")
	}

	if stats.TotalEvents < 2 {
		t.Errorf("Expected >= 2 events, got %d", stats.TotalEvents)
	}
}

func TestInMemoryEventStore_Health(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	health, err := store.Health(ctx)
	if err != nil {
		t.Fatalf("Health() error = %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", health.Status)
	}
}

func TestInMemoryEventStore_CompactEvents(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	ctx := context.Background()

	// Save multiple events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e3", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 3, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	// Compact up to version 2
	err := store.CompactEvents(ctx, "agg-1", 2)
	if err != nil {
		t.Fatalf("CompactEvents() error = %v", err)
	}

	// Verify only version 3 remains
	remaining, err := store.GetAllEvents(ctx, "agg-1")
	if err != nil {
		t.Fatalf("GetAllEvents() error = %v", err)
	}

	if len(remaining) != 1 {
		t.Errorf("Expected 1 remaining event, got %d", len(remaining))
	}

	if remaining[0].Version != 3 {
		t.Errorf("Expected version 3, got %d", remaining[0].Version)
	}
}

func TestDefaultEventStoreConfig(t *testing.T) {
	config := DefaultEventStoreConfig()

	if config == nil {
		t.Fatal("DefaultEventStoreConfig returned nil")
	}

	if config.StoreType == "" {
		t.Error("StoreType should not be empty")
	}
}
