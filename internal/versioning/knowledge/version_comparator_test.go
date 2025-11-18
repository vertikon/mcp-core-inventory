package knowledge

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryVersionComparator_CompareVersions(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	vc := NewInMemoryVersionComparator(kv)
	ctx := context.Background()

	// Create two versions
	v1, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	v2, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Add documents to v1
	doc1 := &KnowledgeDocument{Content: "doc1"}
	err = kv.AddDocument(ctx, v1.ID, doc1)
	require.NoError(t, err)

	// Add documents to v2
	doc2 := &KnowledgeDocument{Content: "doc2"}
	err = kv.AddDocument(ctx, v2.ID, doc2)
	require.NoError(t, err)

	// Compare versions
	diff, err := vc.CompareVersions(ctx, v1.ID, v2.ID)
	require.NoError(t, err)
	assert.NotNil(t, diff)
	assert.Len(t, diff.AddedDocuments, 1)
	assert.Len(t, diff.RemovedDocuments, 1)
}

func TestInMemoryVersionComparator_CompareSemantic(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	vc := NewInMemoryVersionComparator(kv)
	ctx := context.Background()

	// Create two versions with same documents
	v1, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	v2, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	doc := &KnowledgeDocument{Content: "same doc"}
	err = kv.AddDocument(ctx, v1.ID, doc)
	require.NoError(t, err)

	doc2 := &KnowledgeDocument{ID: doc.ID, Content: "same doc"}
	err = kv.AddDocument(ctx, v2.ID, doc2)
	require.NoError(t, err)

	// Compare semantic similarity
	similarity, err := vc.CompareSemantic(ctx, v1.ID, v2.ID)
	require.NoError(t, err)
	assert.Greater(t, similarity, 0.0)
	assert.LessOrEqual(t, similarity, 1.0)
}

func TestInMemoryVersionComparator_CompareStructural(t *testing.T) {
	kv := NewInMemoryKnowledgeVersioning()
	vc := NewInMemoryVersionComparator(kv)
	ctx := context.Background()

	// Create two versions
	v1, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)
	v2, err := kv.CreateVersion(ctx, "kb-1", nil)
	require.NoError(t, err)

	// Add documents
	doc1 := &KnowledgeDocument{Content: "doc1"}
	err = kv.AddDocument(ctx, v1.ID, doc1)
	require.NoError(t, err)

	doc2 := &KnowledgeDocument{Content: "doc2"}
	err = kv.AddDocument(ctx, v2.ID, doc2)
	require.NoError(t, err)

	// Compare structural similarity
	similarity, err := vc.CompareStructural(ctx, v1.ID, v2.ID)
	require.NoError(t, err)
	assert.Equal(t, 1.0, similarity) // Same document count
}
