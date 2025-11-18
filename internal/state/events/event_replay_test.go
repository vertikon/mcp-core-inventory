package events

import (
	"context"
	"testing"
	"time"
)

// MockReplayHandler for testing
type MockReplayHandler struct {
	handledEvents []*Event
}

func (h *MockReplayHandler) Handle(ctx context.Context, event *Event) error {
	h.handledEvents = append(h.handledEvents, event)
	return nil
}

func (h *MockReplayHandler) CanHandle(event *Event) bool {
	return true
}

func (h *MockReplayHandler) GetHandlerType() string {
	return "mock"
}

func TestNewEventReplay(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	replay := NewEventReplay(store, nil)

	if replay == nil {
		t.Fatal("NewEventReplay returned nil")
	}
}

func TestEventReplayImpl_ReplayEvents(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	replay := NewEventReplay(store, nil)
	ctx := context.Background()

	// Save events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	handler := &MockReplayHandler{}

	// Replay events
	progress, err := replay.ReplayEvents(ctx, "agg-1", 1, 2, handler)
	if err != nil {
		t.Fatalf("ReplayEvents() error = %v", err)
	}

	if progress == nil {
		t.Fatal("ReplayEvents() returned nil progress")
	}

	if progress.ProcessedEvents != 2 {
		t.Errorf("Expected 2 processed events, got %d", progress.ProcessedEvents)
	}

	if len(handler.handledEvents) != 2 {
		t.Errorf("Expected 2 handled events, got %d", len(handler.handledEvents))
	}
}

func TestEventReplayImpl_ReplayAllEvents(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	replay := NewEventReplay(store, nil)
	ctx := context.Background()

	// Save events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	handler := &MockReplayHandler{}

	// Replay all events
	progress, err := replay.ReplayAllEvents(ctx, "agg-1", handler)
	if err != nil {
		t.Fatalf("ReplayAllEvents() error = %v", err)
	}

	if progress == nil {
		t.Fatal("ReplayAllEvents() returned nil progress")
	}

	if progress.ProcessedEvents != 2 {
		t.Errorf("Expected 2 processed events, got %d", progress.ProcessedEvents)
	}
}

func TestEventReplayImpl_ReplayFromSnapshot(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	replay := NewEventReplay(store, nil)
	ctx := context.Background()

	// Save events and create snapshot
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
		{ID: "e2", Type: EventTypeUpdate, AggregateID: "agg-1", AggregateType: "test", Version: 2, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)
	_ = store.CreateSnapshot(ctx, "agg-1", 1, map[string]interface{}{"state": "snapshot"})

	// Add more events after snapshot
	_ = store.SaveEvent(ctx, &Event{
		ID:            "e3",
		Type:          EventTypeUpdate,
		AggregateID:   "agg-1",
		AggregateType: "test",
		Version:       3,
		Timestamp:     time.Now(),
		NodeID:        "node-1",
	})

	handler := &MockReplayHandler{}

	// Replay from snapshot
	progress, err := replay.ReplayFromSnapshot(ctx, "agg-1", 1, handler)
	if err != nil {
		t.Fatalf("ReplayFromSnapshot() error = %v", err)
	}

	if progress == nil {
		t.Fatal("ReplayFromSnapshot() returned nil progress")
	}

	// Should only replay events after snapshot (version 3)
	if progress.ProcessedEvents != 1 {
		t.Errorf("Expected 1 processed event, got %d", progress.ProcessedEvents)
	}
}

func TestEventReplayImpl_GetReplayStats(t *testing.T) {
	store := NewInMemoryEventStore(nil)
	replay := NewEventReplay(store, nil)
	ctx := context.Background()

	// Save and replay events
	events := []*Event{
		{ID: "e1", Type: EventTypeCreate, AggregateID: "agg-1", AggregateType: "test", Version: 1, Timestamp: time.Now(), NodeID: "node-1"},
	}
	_ = store.SaveEvents(ctx, events)

	handler := &MockReplayHandler{}
	_, _ = replay.ReplayEvents(ctx, "agg-1", 1, 1, handler)

	// Get stats
	stats, err := replay.GetReplayStats(ctx)
	if err != nil {
		t.Fatalf("GetReplayStats() error = %v", err)
	}

	if stats == nil {
		t.Fatal("GetReplayStats() returned nil")
	}

	if stats.TotalReplays == 0 {
		t.Error("TotalReplays should be > 0")
	}
}

func TestDefaultReplayConfig(t *testing.T) {
	config := DefaultReplayConfig()

	if config == nil {
		t.Fatal("DefaultReplayConfig returned nil")
	}

	if config.Strategy == "" {
		t.Error("Strategy should not be empty")
	}
}
