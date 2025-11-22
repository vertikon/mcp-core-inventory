package validators

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewValidatorFactory(t *testing.T) {
	factory := NewValidatorFactory()
	if factory == nil {
		t.Fatal("NewValidatorFactory returned nil")
	}
	if factory.structureValidator == nil {
		t.Error("structureValidator should not be nil")
	}
	if factory.dependencyValidator == nil {
		t.Error("dependencyValidator should not be nil")
	}
	if factory.treeValidator == nil {
		t.Error("treeValidator should not be nil")
	}
	if factory.securityValidator == nil {
		t.Error("securityValidator should not be nil")
	}
	if factory.configValidator == nil {
		t.Error("configValidator should not be nil")
	}
}

func TestValidatorFactory_GetValidators(t *testing.T) {
	factory := NewValidatorFactory()

	if factory.GetStructureValidator() == nil {
		t.Error("GetStructureValidator() returned nil")
	}
	if factory.GetDependencyValidator() == nil {
		t.Error("GetDependencyValidator() returned nil")
	}
	if factory.GetTreeValidator() == nil {
		t.Error("GetTreeValidator() returned nil")
	}
	if factory.GetSecurityValidator() == nil {
		t.Error("GetSecurityValidator() returned nil")
	}
	if factory.GetConfigValidator() == nil {
		t.Error("GetConfigValidator() returned nil")
	}
}

func TestStructureValidator_Validate(t *testing.T) {
	validator := NewStructureValidator()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "test-project")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some basic structure
	files := []string{
		"go.mod",
		"README.md",
		".gitignore",
		"Makefile",
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, file), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}
	dirs := []string{
		"cmd",
		filepath.Join("internal", "domain"),
		filepath.Join("internal", "application"),
		filepath.Join("internal", "infrastructure"),
		"configs",
		"pkg",
		"tests",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name      string
		request   StructureRequest
		wantErr   bool
		wantValid bool
	}{
		{
			name: "valid structure",
			request: StructureRequest{
				Path:       tmpDir,
				StrictMode: false,
			},
			wantErr:   false,
			wantValid: true,
		},
		{
			name: "nonexistent path",
			request: StructureRequest{
				Path:       "/nonexistent/path",
				StrictMode: false,
			},
			wantErr:   false,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result == nil {
				t.Fatal("Validate() returned nil result")
			}
			if result.Valid != tt.wantValid {
				t.Errorf("Validate() valid = %v, want %v (errors=%v, warnings=%v)", result.Valid, tt.wantValid, result.Errors, result.Warnings)
			}
		})
	}
}

func TestDependencyValidator_Validate(t *testing.T) {
	validator := NewDependencyValidator()

	tmpDir, err := os.MkdirTemp("", "test-project")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	request := DependencyRequest{
		Path: tmpDir,
	}

	result, err := validator.Validate(context.Background(), request)
	if err != nil {
		t.Logf("DependencyValidator.Validate() error = %v (may be expected)", err)
	}
	if result == nil {
		t.Error("Validate() should return result")
	}
}

func TestValidationResult(t *testing.T) {
	result := ValidationResult{
		Path:        "/test/path",
		Valid:       true,
		Warnings:    []string{"warning1"},
		Errors:      []string{},
		Checks:      []string{"structure"},
		Duration:    time.Second,
		ValidatedAt: time.Now(),
	}

	if result.Path != "/test/path" {
		t.Errorf("Expected path '/test/path', got '%s'", result.Path)
	}
	if !result.Valid {
		t.Error("Expected Valid to be true")
	}
	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
	}
}
