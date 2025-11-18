package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryModelRegistry_RegisterModel(t *testing.T) {
	tests := []struct {
		name      string
		model     *Model
		wantError bool
	}{
		{
			name: "successful registration",
			model: &Model{
				Name:        "test-model",
				Type:        "classification",
				Description: "test model",
			},
			wantError: false,
		},
		{
			name: "registration with metadata",
			model: &Model{
				Name:     "test-model-2",
				Type:     "regression",
				Metadata: map[string]interface{}{"framework": "pytorch"},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mr := NewInMemoryModelRegistry()
			ctx := context.Background()

			model, err := mr.RegisterModel(ctx, tt.model)
			if tt.wantError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, model.ID)
			assert.Equal(t, tt.model.Name, model.Name)
		})
	}
}

func TestInMemoryModelRegistry_RegisterVersion(t *testing.T) {
	mr := NewInMemoryModelRegistry()
	ctx := context.Background()

	// Register model
	model, err := mr.RegisterModel(ctx, &Model{
		Name: "test-model",
		Type: "classification",
	})
	require.NoError(t, err)

	// Register version
	version, err := mr.RegisterVersion(ctx, model.ID, &ModelVersion{
		Path: "/path/to/model",
	})
	require.NoError(t, err)
	assert.NotEmpty(t, version.ID)
	assert.Equal(t, model.ID, version.ModelID)
	assert.Equal(t, "v1", version.Version)
}

func TestInMemoryModelRegistry_GetLatestVersion(t *testing.T) {
	mr := NewInMemoryModelRegistry()
	ctx := context.Background()

	// Register model
	model, err := mr.RegisterModel(ctx, &Model{
		Name: "test-model",
		Type: "classification",
	})
	require.NoError(t, err)

	// Register multiple versions
	_, err = mr.RegisterVersion(ctx, model.ID, &ModelVersion{Path: "/v1"})
	require.NoError(t, err)
	_, err = mr.RegisterVersion(ctx, model.ID, &ModelVersion{Path: "/v2"})
	require.NoError(t, err)

	// Get latest
	latest, err := mr.GetLatestVersion(ctx, model.ID)
	require.NoError(t, err)
	assert.Equal(t, "v2", latest.Version)
}

func TestInMemoryModelRegistry_CalculateFingerprint(t *testing.T) {
	mr := NewInMemoryModelRegistry()

	data := []byte("test data")
	fingerprint := mr.CalculateFingerprint(data)
	assert.NotEmpty(t, fingerprint)
	assert.Len(t, fingerprint, 64) // SHA256 hex string
}
