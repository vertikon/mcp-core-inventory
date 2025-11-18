// Package value_objects provides domain value objects
package value_objects

import (
	"fmt"
	"time"
)

// FeatureStatus represents the status of a feature
type FeatureStatus string

const (
	FeatureStatusEnabled  FeatureStatus = "enabled"
	FeatureStatusDisabled FeatureStatus = "disabled"
)

// Feature represents a project feature configuration
type Feature struct {
	name        string
	status      FeatureStatus
	config      map[string]interface{}
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewFeature creates a new feature
func NewFeature(name string, status FeatureStatus, description string) (*Feature, error) {
	if name == "" {
		return nil, fmt.Errorf("feature name cannot be empty")
	}

	return &Feature{
		name:        name,
		status:      status,
		config:      make(map[string]interface{}),
		description: description,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}, nil
}

// Name returns the feature name
func (f *Feature) Name() string {
	return f.name
}

// Status returns the feature status
func (f *Feature) Status() FeatureStatus {
	return f.status
}

// Description returns the feature description
func (f *Feature) Description() string {
	return f.description
}

// Config returns the feature configuration
func (f *Feature) Config() map[string]interface{} {
	// Return a copy to maintain immutability
	config := make(map[string]interface{})
	for k, v := range f.config {
		config[k] = v
	}
	return config
}

// SetConfig sets a configuration value
func (f *Feature) SetConfig(key string, value interface{}) error {
	if key == "" {
		return fmt.Errorf("config key cannot be empty")
	}
	f.config[key] = value
	f.updatedAt = time.Now()
	return nil
}

// Enable enables the feature
func (f *Feature) Enable() {
	f.status = FeatureStatusEnabled
	f.updatedAt = time.Now()
}

// Disable disables the feature
func (f *Feature) Disable() {
	f.status = FeatureStatusDisabled
	f.updatedAt = time.Now()
}

// IsEnabled checks if the feature is enabled
func (f *Feature) IsEnabled() bool {
	return f.status == FeatureStatusEnabled
}

// CreatedAt returns the creation timestamp
func (f *Feature) CreatedAt() time.Time {
	return f.createdAt
}

// UpdatedAt returns the last update timestamp
func (f *Feature) UpdatedAt() time.Time {
	return f.updatedAt
}

// Equals checks if two features are equal
func (f *Feature) Equals(other *Feature) bool {
	if other == nil {
		return false
	}
	return f.name == other.name
}
