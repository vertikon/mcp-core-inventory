package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewCoherencyManager(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)

	if manager == nil {
		t.Fatal("NewCoherencyManager returned nil")
	}
}

func TestCoherencyManagerImpl_Invalidate(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Set value in cache
	_ = cache.Set(ctx, "key1", "value1", CacheLevelL1, 5*time.Minute)

	// Invalidate
	err := manager.Invalidate(ctx, "key1", "test")
	if err != nil {
		t.Fatalf("Invalidate() error = %v", err)
	}

	// Verify key is invalidated
	_, err = cache.Get(ctx, "key1")
	if err == nil {
		t.Error("Get() should fail for invalidated key")
	}
}

func TestCoherencyManagerImpl_InvalidateAll(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Set values
	_ = cache.Set(ctx, "key1", "value1", CacheLevelL1, 5*time.Minute)
	_ = cache.Set(ctx, "key2", "value2", CacheLevelL1, 5*time.Minute)

	// Invalidate all
	err := manager.InvalidateAll(ctx)
	if err != nil {
		t.Fatalf("InvalidateAll() error = %v", err)
	}
}

func TestCoherencyManagerImpl_Update(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Update cache
	err := manager.Update(ctx, "key1", "updated-value")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func TestCoherencyManagerImpl_GetCoherencyStatus(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	status, err := manager.GetCoherencyStatus(ctx)
	if err != nil {
		t.Fatalf("GetCoherencyStatus() error = %v", err)
	}

	if status == nil {
		t.Fatal("GetCoherencyStatus() returned nil")
	}

	if status.Strategy == "" {
		t.Error("Strategy should not be empty")
	}
}

func TestCoherencyManagerImpl_GetInvalidationStats(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Perform invalidations
	_ = manager.Invalidate(ctx, "key1", "test1")
	_ = manager.Invalidate(ctx, "key2", "test2")

	// Get stats
	stats, err := manager.GetInvalidationStats(ctx)
	if err != nil {
		t.Fatalf("GetInvalidationStats() error = %v", err)
	}

	if stats == nil {
		t.Fatal("GetInvalidationStats() returned nil")
	}

	if stats.TotalInvalidations < 2 {
		t.Errorf("Expected >= 2 invalidations, got %d", stats.TotalInvalidations)
	}
}

func TestCoherencyManagerImpl_OnStoreUpdate(t *testing.T) {
	config := DefaultCoherencyConfig()
	config.Strategy = CoherencyStrategyInvalidate
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Set value in cache
	_ = cache.Set(ctx, "key1", "value1", CacheLevelL1, 5*time.Minute)

	// Simulate store update
	err := manager.OnStoreUpdate(ctx, "key1", "new-value")
	if err != nil {
		t.Fatalf("OnStoreUpdate() error = %v", err)
	}
}

func TestCoherencyManagerImpl_StartBackgroundInvalidator_StopBackgroundInvalidator(t *testing.T) {
	config := DefaultCoherencyConfig()
	cache := NewStateCache(nil)
	manager := NewCoherencyManager(config, cache)
	ctx := context.Background()

	// Start background invalidator
	err := manager.StartBackgroundInvalidator(ctx)
	if err != nil {
		t.Fatalf("StartBackgroundInvalidator() error = %v", err)
	}

	// Stop background invalidator
	err = manager.StopBackgroundInvalidator()
	if err != nil {
		t.Fatalf("StopBackgroundInvalidator() error = %v", err)
	}
}

func TestDefaultCoherencyConfig(t *testing.T) {
	config := DefaultCoherencyConfig()

	if config == nil {
		t.Fatal("DefaultCoherencyConfig returned nil")
	}

	if config.Strategy == "" {
		t.Error("Strategy should not be empty")
	}
}
