package knowledge

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryKnowledgeVersioning_CreateVersion(t *testing.T) {
	tests := []struct {
		name        string
		knowledgeID string
		metadata    map[string]interface{}
		wantError   bool
	}{
		{
			name:        "successful creation",
			knowledgeID: "kb-1",
			metadata:    map[string]interface{}{"description": "test"},
			wantError:   false,
		},
		{
			name:        "creation without metadata",
			knowledgeID: "kb-2",
			metadata:    nil,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := NewInMemoryKnowledgeVersioning()
			ctx := context.Background()

			version, err := kv.CreateVersion(ctx, tt.knowledgeID, tt.metadata)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, version.ID)
			assert.Equal(t, tt.knowledgeID, version.KnowledgeID)
			assert.Equal(t, "v1", version.Version)
			assert.NotZero(t, version.CreatedAt)
		})
	}
}

func TestInMemoryKnowledgeVersioning_GetVersion(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	ctx := context.Background()

	// Create a version
	version, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Get the version
	got, err := kv.GetVersion(ctx, version.ID)
	require.NoError(t, err)
	assert.Equal(t, version.ID, got.ID)
	assert.Equal(t, version.KnowledgeID, got.KnowledgeID)

	// Try to get non-existent version
	_, err = kv.GetVersion(ctx, "non-existent")
	assert.Error(t, err)
}

func TestInMemoryKnowledgeVersioning_AddDocument(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	ctx := context.Background()

	// Create a version
	version, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Add document
	doc := &KnowledgeDocument{
		Content:  "test content",
		Metadata: map[string]interface{}{"source": "test"},
	}

	err = kv.AddDocument(ctx, version.ID, doc)
	require.NoError(t, err)
	assert.NotEmpty(t, doc.ID)

	// Verify document count
	got, err := kv.GetVersion(ctx, version.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, got.DocumentCount)
}

func TestInMemoryKnowledgeVersioning_ListVersions(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	ctx := context.Background()

	// Create multiple versions
	_, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	_, err = kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	_, err = kv.CreateVersion(ctx, "kb-2", nil)
	require.NoError(t, err)

	// List versions for kb-1
	versions, err := kv.ListVersions(ctx, "kb-1")
	require.NoError(t, err)
	assert.Len(t, versions, 2)

	// List versions for kb-2
	versions, err = kv.ListVersions(ctx, "kb-2")
	require.NoError(t, err)
	assert.Len(t, versions, 1)
}

func TestInMemoryKnowledgeVersioning_GetLatestVersion(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	ctx := context.Background()

	// Create multiple versions
	_, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	v2, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Get latest
	latest, err := kv.GetLatestVersion(ctx, "kb-1")
	require.NoError(t, err)
	assert.Equal(t, v2.ID, latest.ID)
	assert.Equal(t, "v2", latest.Version)
}

func TestInMemoryKnowledgeVersioning_DeleteVersion(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	ctx := context.Background()

	// Create a version
	version, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Delete version
	err = kv.DeleteVersion(ctx, version.ID)
	require.NoError(t, err)

	// Verify soft delete
	got, err := kv.GetVersion(ctx, version.ID)
	require.NoError(t, err)
	assert.True(t, got.Metadata["deleted"].(bool))
}
