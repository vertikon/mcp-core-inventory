// Package config provides configuration loading and validation
package config

import (
	"os"
	"strings"
)

// EnvironmentManager manages environment-specific configuration
type EnvironmentManager struct {
	env string
}

// NewEnvironmentManager creates a new environment manager
func NewEnvironmentManager() *EnvironmentManager {
	env := os.Getenv("HULK_ENV")
	if env == "" {
		env = os.Getenv("MCP_HULK_ENV") // fallback for backward compatibility
	}
	if env == "" {
		env = "development"
	}

	return &EnvironmentManager{env: strings.ToLower(env)}
}

// GetEnvironment returns the current environment
func (em *EnvironmentManager) GetEnvironment() string {
	return em.env
}

// IsDevelopment returns true if in development mode
func (em *EnvironmentManager) IsDevelopment() bool {
	return em.env == "development" || em.env == "dev"
}

// IsProduction returns true if in production mode
func (em *EnvironmentManager) IsProduction() bool {
	return em.env == "production" || em.env == "prod"
}

// IsStaging returns true if in staging mode
func (em *EnvironmentManager) IsStaging() bool {
	return em.env == "staging" || em.env == "stage"
}

// IsTest returns true if in test mode
func (em *EnvironmentManager) IsTest() bool {
	return em.env == "test"
}
