package data

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryDataVersioning_CreateVersion(t *testing.T) {
	dv := NewInMemoryDataVersioning()
	ctx := context.Background()

	version, err := dv.CreateVersion(ctx, "dataset-1", map[string]interface{}{"description": "test"})
	require.NoError(t, err)
	assert.NotEmpty(t, version.ID)
	assert.Equal(t, "dataset-1", version.DatasetID)
	assert.Equal(t, "v1", version.Version)
}

func TestInMemoryDataVersioning_CreateSnapshot(t *testing.T) {
	dv := NewInMemoryDataVersioning()
	ctx := context.Background()

	// Create version
	version, err := dv.CreateVersion(ctx, "dataset-1", nil)
	require.NoError(t, err)

	// Create snapshot
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	snapshot, err := dv.CreateSnapshot(ctx, version.ID, SnapshotTypeFull, data)
	require.NoError(t, err)
	assert.NotEmpty(t, snapshot.ID)
	assert.Equal(t, version.ID, snapshot.VersionID)
	assert.NotEmpty(t, snapshot.Checksum)
}

func TestInMemoryDataVersioning_GetLatestVersion(t *testing.T) {
	dv := NewInMemoryDataVersioning()
	ctx := context.Background()

	// Create multiple versions
	_, err := dv.CreateVersion(ctx, "dataset-1", nil)
	require.NoError(t, err)
	v2, err := dv.CreateVersion(ctx, "dataset-1", nil)
	require.NoError(t, err)

	// Get latest
	latest, err := dv.GetLatestVersion(ctx, "dataset-1")
	require.NoError(t, err)
	assert.Equal(t, v2.ID, latest.ID)
	assert.Equal(t, "v2", latest.Version)
}

func TestInMemoryDataVersioning_TagVersion(t *testing.T) {
	dv := NewInMemoryDataVersioning()
	ctx := context.Background()

	// Create version
	version, err := dv.CreateVersion(ctx, "dataset-1", nil)
	require.NoError(t, err)

	// Tag version
	err = dv.TagVersion(ctx, version.ID, []string{"production", "stable"})
	require.NoError(t, err)

	// Verify tags
	got, err := dv.GetVersion(ctx, version.ID)
	require.NoError(t, err)
	assert.Contains(t, got.Tags, "production")
	assert.Contains(t, got.Tags, "stable")
}
