package core

import (
	"context"
	"testing"
)

func TestNewRouter(t *testing.T) {
	configs := make(map[LLMProvider]*ProviderConfig)
	strategy := StrategyBalanced
	metrics := NewMetrics()

	router := NewRouter(configs, strategy, metrics)

	if router == nil {
		t.Fatal("NewRouter returned nil")
	}
	if router.configs == nil {
		t.Error("configs map is nil")
	}
	if router.strategy != strategy {
		t.Errorf("Expected strategy %v, got %v", strategy, router.strategy)
	}
	if router.metrics == nil {
		t.Error("metrics is nil")
	}
}

func TestRouter_SelectProvider_ByPriority(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider:     ProviderOpenAI,
			Models:       []string{"gpt-4"},
			DefaultModel: "gpt-4",
			Priority:     10,
			Enabled:      true,
		},
		ProviderGemini: {
			Provider:     ProviderGemini,
			Models:       []string{"gemini-pro"},
			DefaultModel: "gemini-pro",
			Priority:     5,
			Enabled:      true,
		},
	}
	router := NewRouter(configs, StrategyBalanced, NewMetrics())
	router.SetAvailability(ProviderOpenAI, true)
	router.SetAvailability(ProviderGemini, true)

	req := &LLMRequest{
		Prompt: "test",
	}

	provider, model, err := router.SelectProvider(context.Background(), req)
	if err != nil {
		t.Fatalf("SelectProvider failed: %v", err)
	}

	if provider != ProviderOpenAI {
		t.Errorf("Expected %v (higher priority), got %v", ProviderOpenAI, provider)
	}
	if model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", model)
	}
}

func TestRouter_SelectProvider_ByCost(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider:     ProviderOpenAI,
			Models:       []string{"gpt-4"},
			DefaultModel: "gpt-4",
			CostPerToken: 0.03,
			Priority:     10,
			Enabled:      true,
		},
		ProviderGemini: {
			Provider:     ProviderGemini,
			Models:       []string{"gemini-pro"},
			DefaultModel: "gemini-pro",
			CostPerToken: 0.01,
			Priority:     5,
			Enabled:      true,
		},
	}
	router := NewRouter(configs, StrategyCost, NewMetrics())
	router.SetAvailability(ProviderOpenAI, true)
	router.SetAvailability(ProviderGemini, true)

	req := &LLMRequest{
		Prompt: "test",
	}

	provider, _, err := router.SelectProvider(context.Background(), req)
	if err != nil {
		t.Fatalf("SelectProvider failed: %v", err)
	}

	if provider != ProviderGemini {
		t.Errorf("Expected %v (lower cost), got %v", ProviderGemini, provider)
	}
}

func TestRouter_SelectProvider_NoAvailable(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider: ProviderOpenAI,
			Models:   []string{"gpt-4"},
			Enabled:  true,
		},
	}
	router := NewRouter(configs, StrategyBalanced, NewMetrics())
	router.SetAvailability(ProviderOpenAI, false)

	req := &LLMRequest{
		Prompt: "test",
	}

	_, _, err := router.SelectProvider(context.Background(), req)
	if err == nil {
		t.Error("Expected error when no providers available")
	}
}

func TestRouter_SelectFallback(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider:      ProviderOpenAI,
			Models:        []string{"gpt-4"},
			DefaultModel:  "gpt-4",
			Priority:      10,
			Enabled:       true,
			FallbackOrder: 1,
		},
		ProviderGemini: {
			Provider:      ProviderGemini,
			Models:        []string{"gemini-pro"},
			DefaultModel:  "gemini-pro",
			Priority:      5,
			Enabled:       true,
			FallbackOrder: 0,
		},
	}
	router := NewRouter(configs, StrategyBalanced, NewMetrics())
	router.SetAvailability(ProviderOpenAI, false)
	router.SetAvailability(ProviderGemini, true)

	req := &LLMRequest{
		Prompt: "test",
	}

	provider, model, err := router.SelectFallback(context.Background(), req, ProviderOpenAI)
	if err != nil {
		t.Fatalf("SelectFallback failed: %v", err)
	}

	if provider != ProviderGemini {
		t.Errorf("Expected fallback to %v, got %v", ProviderGemini, provider)
	}
	if model != "gemini-pro" {
		t.Errorf("Expected model 'gemini-pro', got '%s'", model)
	}
}

func TestRouter_SelectProvider_WithRequestedModel(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider:     ProviderOpenAI,
			Models:       []string{"gpt-4", "gpt-3.5-turbo"},
			DefaultModel: "gpt-4",
			Priority:     10,
			Enabled:      true,
		},
	}
	router := NewRouter(configs, StrategyBalanced, NewMetrics())
	router.SetAvailability(ProviderOpenAI, true)

	req := &LLMRequest{
		Prompt: "test",
		Model:  "gpt-3.5-turbo",
	}

	_, model, err := router.SelectProvider(context.Background(), req)
	if err != nil {
		t.Fatalf("SelectProvider failed: %v", err)
	}

	if model != "gpt-3.5-turbo" {
		t.Errorf("Expected requested model 'gpt-3.5-turbo', got '%s'", model)
	}
}

func TestRouter_SetAvailability(t *testing.T) {
	configs := map[LLMProvider]*ProviderConfig{
		ProviderOpenAI: {
			Provider: ProviderOpenAI,
			Enabled:  true,
		},
	}
	router := NewRouter(configs, StrategyBalanced, NewMetrics())

	router.SetAvailability(ProviderOpenAI, true)
	if !router.availability[ProviderOpenAI] {
		t.Error("Availability not set to true")
	}

	router.SetAvailability(ProviderOpenAI, false)
	if router.availability[ProviderOpenAI] {
		t.Error("Availability not set to false")
	}
}
