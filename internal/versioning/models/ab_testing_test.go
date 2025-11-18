package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryABTesting_CreateTest(t *testing.T) {
	ab := NewInMemoryABTesting()
	ctx := context.Background()

	test := &ABTest{
		Name:     "test-ab",
		ModelID:  "model-1",
		VersionA: "version-a",
		VersionB: "version-b",
		TrafficSplit: TrafficSplit{
			VersionAPercent: 0.5,
			VersionBPercent: 0.5,
		},
		Criteria: PromotionCriteria{
			MinRequests:    100,
			MinScore:       0.8,
			MaxErrorRate:   0.05,
			MaxLatencyMs:   1000,
			MinImprovement: 5.0,
		},
	}

	created, err := ab.CreateTest(ctx, test)
	require.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, ABTestStatusDraft, created.Status)
}

func TestInMemoryABTesting_StartTest(t *testing.T) {
	ab := NewInMemoryABTesting()
	ctx := context.Background()

	test := &ABTest{
		Name:     "test-ab",
		ModelID:  "model-1",
		VersionA: "version-a",
		VersionB: "version-b",
		TrafficSplit: TrafficSplit{
			VersionAPercent: 0.5,
			VersionBPercent: 0.5,
		},
	}
	created, err := ab.CreateTest(ctx, test)
	require.NoError(t, err)

	err = ab.StartTest(ctx, created.ID)
	require.NoError(t, err)

	got, err := ab.GetTest(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, ABTestStatusRunning, got.Status)
}

func TestInMemoryABTesting_RecordRequest(t *testing.T) {
	ab := NewInMemoryABTesting()
	ctx := context.Background()

	test := &ABTest{
		Name:     "test-ab",
		ModelID:  "model-1",
		VersionA: "version-a",
		VersionB: "version-b",
		TrafficSplit: TrafficSplit{
			VersionAPercent: 0.5,
			VersionBPercent: 0.5,
		},
	}
	created, err := ab.CreateTest(ctx, test)
	require.NoError(t, err)

	err = ab.StartTest(ctx, created.ID)
	require.NoError(t, err)

	err = ab.RecordRequest(ctx, created.ID, "version-a", 100*time.Millisecond, false)
	require.NoError(t, err)

	metrics, err := ab.GetMetrics(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), metrics.VersionARequests)
}

func TestInMemoryABTesting_SelectVersion(t *testing.T) {
	ab := NewInMemoryABTesting()
	ctx := context.Background()

	test := &ABTest{
		Name:     "test-ab",
		ModelID:  "model-1",
		VersionA: "version-a",
		VersionB: "version-b",
		TrafficSplit: TrafficSplit{
			VersionAPercent: 0.5,
			VersionBPercent: 0.5,
		},
	}
	created, err := ab.CreateTest(ctx, test)
	require.NoError(t, err)

	err = ab.StartTest(ctx, created.ID)
	require.NoError(t, err)

	// Select version multiple times to test distribution
	selectedA := 0
	selectedB := 0
	for i := 0; i < 100; i++ {
		version, err := ab.SelectVersion(ctx, created.ID)
		require.NoError(t, err)
		if version == "version-a" {
			selectedA++
		} else {
			selectedB++
		}
	}

	// Should have some distribution (not deterministic but should have both)
	assert.Greater(t, selectedA, 0)
	assert.Greater(t, selectedB, 0)
}
