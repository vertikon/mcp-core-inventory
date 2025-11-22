// Package value_objects provides value objects tests
package value_objects

import (
	"testing"
	"time"
)

func TestNewFeature(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		status      FeatureStatus
		description string
		wantErr     bool
	}{
		{"Valid feature", "feature-valid", FeatureStatusEnabled, "Test feature", false},
		{"Empty name", "", FeatureStatusEnabled, "Test feature", true},
		{"Disabled feature", "feature-disabled", FeatureStatusDisabled, "Test feature", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feature, err := NewFeature(tt.inputName, tt.status, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFeature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && feature == nil {
				t.Errorf("NewFeature() returned nil feature")
			}
			if !tt.wantErr && feature.Name() != tt.inputName {
				t.Errorf("NewFeature() name = %v, want %v", feature.Name(), tt.inputName)
			}
		})
	}
}

func TestFeature_EnableDisable(t *testing.T) {
	feature, _ := NewFeature("test", FeatureStatusDisabled, "Test")

	feature.Enable()
	if !feature.IsEnabled() {
		t.Errorf("Enable() feature not enabled")
	}

	feature.Disable()
	if feature.IsEnabled() {
		t.Errorf("Disable() feature still enabled")
	}
}

func TestFeature_SetConfig(t *testing.T) {
	feature, _ := NewFeature("test", FeatureStatusEnabled, "Test")

	if err := feature.SetConfig("key1", "value1"); err != nil {
		t.Errorf("SetConfig() error = %v", err)
	}

	config := feature.Config()
	if config["key1"] != "value1" {
		t.Errorf("SetConfig() config value = %v, want value1", config["key1"])
	}

	// Test empty key
	if err := feature.SetConfig("", "value"); err == nil {
		t.Errorf("SetConfig() should fail on empty key")
	}
}

func TestFeature_UpdatedAt(t *testing.T) {
	feature, _ := NewFeature("test", FeatureStatusEnabled, "Test")
	initialTime := feature.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	feature.Enable()

	if !feature.UpdatedAt().After(initialTime) {
		t.Errorf("UpdatedAt() should be updated after Enable()")
	}
}
