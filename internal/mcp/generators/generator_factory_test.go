package generators

import (
	"testing"
)

func TestNewGeneratorFactory(t *testing.T) {
	tests := []struct {
		name   string
		config *FactoryConfig
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
		},
		{
			name: "custom config",
			config: &FactoryConfig{
				TemplateRoot:  "./custom-templates",
				DefaultStack:  "web",
				CacheEnabled:  false,
				MaxConcurrent: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewGeneratorFactory(tt.config)
			if factory == nil {
				t.Fatal("NewGeneratorFactory returned nil")
			}
			if factory.generators == nil {
				t.Error("generators map should not be nil")
			}
		})
	}
}

func TestGeneratorFactory_RegisterGenerator(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	generator := NewGoGenerator("./templates/go")
	if err := factory.RegisterGenerator("go", generator); err != nil {
		t.Fatalf("RegisterGenerator() error = %v", err)
	}

	// Try to register nil generator (should fail)
	if err := factory.RegisterGenerator("nil", nil); err == nil {
		t.Error("RegisterGenerator() should fail for nil generator")
	}
}

func TestGeneratorFactory_GetGenerator(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	// Get non-existent generator
	_, err := factory.GetGenerator("nonexistent")
	if err == nil {
		t.Error("GetGenerator() should fail for nonexistent generator")
	}

	// Register and get generator
	generator := NewGoGenerator("./templates/go")
	_ = factory.RegisterGenerator("go", generator)

	got, err := factory.GetGenerator("go")
	if err != nil {
		t.Fatalf("GetGenerator() error = %v", err)
	}
	if got == nil {
		t.Error("GetGenerator() returned nil")
	}
}

func TestGeneratorFactory_ListGenerators(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	// Should have default generators
	generators := factory.ListGenerators()
	if len(generators) == 0 {
		t.Error("ListGenerators() should return default generators")
	}

	// Check for expected default generators
	expected := []string{"go", "web", "tinygo", "wasm", "mcp-go-premium"}
	generatorMap := make(map[string]bool)
	for _, g := range generators {
		generatorMap[g] = true
	}

	for _, exp := range expected {
		if !generatorMap[exp] {
			t.Errorf("Expected generator '%s' not found", exp)
		}
	}
}

func TestGeneratorFactory_GetGeneratorInfo(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	info, err := factory.GetGeneratorInfo("go")
	if err != nil {
		t.Fatalf("GetGeneratorInfo() error = %v", err)
	}

	if info == nil {
		t.Fatal("GetGeneratorInfo() returned nil")
	}

	if stack, ok := info["stack"].(string); !ok || stack != "go" {
		t.Errorf("Expected stack 'go', got '%v'", info["stack"])
	}
}

func TestGeneratorFactory_ValidateRequest(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	tests := []struct {
		name    string
		request GenerateRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: GenerateRequest{
				Name:  "test-project",
				Stack: "go",
				Path:  "/tmp",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: GenerateRequest{
				Stack: "go",
				Path:  "/tmp",
			},
			wantErr: true,
		},
		{
			name: "missing path",
			request: GenerateRequest{
				Name:  "test-project",
				Stack: "go",
			},
			wantErr: true,
		},
		{
			name: "unsupported stack",
			request: GenerateRequest{
				Name:  "test-project",
				Stack: "nonexistent",
				Path:  "/tmp",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := factory.ValidateRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGeneratorFactory_HasGenerator(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	if !factory.HasGenerator("go") {
		t.Error("HasGenerator() should return true for 'go'")
	}

	if factory.HasGenerator("nonexistent") {
		t.Error("HasGenerator() should return false for nonexistent generator")
	}
}

func TestGeneratorFactory_GetDefaultStack(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	defaultStack := factory.GetDefaultStack()
	if defaultStack == "" {
		t.Error("GetDefaultStack() should return a default stack")
	}
}

func TestGeneratorFactory_SetDefaultStack(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	// Set valid stack
	if err := factory.SetDefaultStack("go"); err != nil {
		t.Errorf("SetDefaultStack() error = %v", err)
	}

	// Set invalid stack
	if err := factory.SetDefaultStack("nonexistent"); err == nil {
		t.Error("SetDefaultStack() should fail for nonexistent stack")
	}
}

func TestGeneratorFactory_GetFactoryStats(t *testing.T) {
	factory := NewGeneratorFactory(nil)

	stats := factory.GetFactoryStats()
	if stats.TotalGenerators == 0 {
		t.Error("GetFactoryStats() should return generators count > 0")
	}
	if len(stats.AvailableStacks) == 0 {
		t.Error("GetFactoryStats() should return available stacks")
	}
}
