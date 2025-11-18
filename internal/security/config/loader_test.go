package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: time package is imported for ParseDuration tests

func TestLoader_LoadAuthConfig(t *testing.T) {
	loader := NewLoader()

	cfg, err := loader.LoadAuthConfig()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.JWT.SecretKey)
	assert.Equal(t, "HS256", cfg.JWT.SigningMethod)
	assert.Equal(t, "1h", cfg.JWT.TokenTTL)
	assert.Equal(t, "24h", cfg.JWT.RefreshTTL)
	assert.Equal(t, "24h", cfg.Session.TTL)
	assert.Equal(t, 5, cfg.Session.MaxSessionsPerUser)
}

func TestLoader_LoadRBACConfig(t *testing.T) {
	loader := NewLoader()

	cfg, err := loader.LoadRBACConfig()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.GreaterOrEqual(t, len(cfg.Roles), 2) // Should have at least admin and user roles
}

func TestLoader_LoadEncryptionConfig(t *testing.T) {
	loader := NewLoader()

	cfg, err := loader.LoadEncryptionConfig()

	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "AES-256-GCM", cfg.Algorithm)
	assert.Equal(t, 2048, cfg.RSAKeySize)
}

func TestLoader_ResolveEnvVars(t *testing.T) {
	loader := NewLoader()

	// Set test environment variable
	os.Setenv("MCP_HULK_JWT_SECRET_KEY", "test-secret-from-env")
	defer os.Unsetenv("MCP_HULK_JWT_SECRET_KEY")

	cfg, err := loader.LoadAuthConfig()

	require.NoError(t, err)
	// Should use environment variable if set
	assert.NotNil(t, cfg)
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "valid hours",
			input:    "24h",
			expected: 24 * time.Hour,
			wantErr:  false,
		},
		{
			name:     "valid minutes",
			input:    "30m",
			expected: 30 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "invalid format",
			input:    "invalid",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := ParseDuration(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, duration)
			}
		})
	}
}
