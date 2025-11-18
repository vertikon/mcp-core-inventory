package core

import (
	"errors"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()

	if m == nil {
		t.Fatal("NewMetrics returned nil")
	}
	if m.generations == nil {
		t.Error("generations map is nil")
	}
	if m.totalTokens == nil {
		t.Error("totalTokens map is nil")
	}
}

func TestMetrics_RecordGeneration(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 200*time.Millisecond, true)

	key := m.key(ProviderOpenAI, "gpt-4")
	if m.generations[key] != 1 {
		t.Errorf("Expected 1 generation, got %d", m.generations[key])
	}
	if m.totalTokens[key] != 100 {
		t.Errorf("Expected 100 tokens, got %d", m.totalTokens[key])
	}
	if m.successCount[key] != 1 {
		t.Errorf("Expected 1 success, got %d", m.successCount[key])
	}
}

func TestMetrics_RecordError(t *testing.T) {
	m := NewMetrics()

	err := errors.New("test error")
	m.RecordError(ProviderOpenAI, "gpt-4", err)

	key := m.key(ProviderOpenAI, "gpt-4")
	if m.errorCount[key] != 1 {
		t.Errorf("Expected 1 error, got %d", m.errorCount[key])
	}
}

func TestMetrics_GetTotalGenerations(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 100*time.Millisecond, true)
	m.RecordGeneration(ProviderOpenAI, "gpt-4", 50, 50*time.Millisecond, true)

	total := m.GetTotalGenerations(ProviderOpenAI, "gpt-4")
	if total != 2 {
		t.Errorf("Expected 2 generations, got %d", total)
	}
}

func TestMetrics_GetTotalTokens(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 100*time.Millisecond, true)
	m.RecordGeneration(ProviderOpenAI, "gpt-4", 50, 50*time.Millisecond, true)

	total := m.GetTotalTokens(ProviderOpenAI, "gpt-4")
	if total != 150 {
		t.Errorf("Expected 150 tokens, got %d", total)
	}
}

func TestMetrics_GetAverageLatency(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 200*time.Millisecond, true)
	m.RecordGeneration(ProviderOpenAI, "gpt-4", 50, 100*time.Millisecond, true)

	avg := m.GetAverageLatency(ProviderOpenAI, "gpt-4")
	expected := 150 * time.Millisecond
	if avg != expected {
		t.Errorf("Expected average latency %v, got %v", expected, avg)
	}
}

func TestMetrics_GetSuccessRate(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 100*time.Millisecond, true)
	m.RecordGeneration(ProviderOpenAI, "gpt-4", 50, 50*time.Millisecond, false)
	m.RecordError(ProviderOpenAI, "gpt-4", errors.New("error"))

	rate := m.GetSuccessRate(ProviderOpenAI, "gpt-4")
	// successCount = 1, errorCount = 2 (1 from RecordGeneration with success=false + 1 from RecordError)
	// total = 3, success rate = 1/3 = 0.333...
	expected := 1.0 / 3.0
	if rate != expected {
		t.Errorf("Expected success rate %f, got %f", expected, rate)
	}
}

func TestMetrics_GetErrorRate(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 100*time.Millisecond, true)
	m.RecordError(ProviderOpenAI, "gpt-4", errors.New("error"))

	rate := m.GetErrorRate(ProviderOpenAI, "gpt-4")
	expected := 0.5 // 1 error / 2 total
	if rate != expected {
		t.Errorf("Expected error rate %f, got %f", expected, rate)
	}
}

func TestMetrics_GetRecentErrors(t *testing.T) {
	m := NewMetrics()

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	m.RecordError(ProviderOpenAI, "gpt-4", err1)
	m.RecordError(ProviderOpenAI, "gpt-4", err2)
	m.RecordError(ProviderOpenAI, "gpt-4", err3)

	recent := m.GetRecentErrors(ProviderOpenAI, "gpt-4", 2)
	if len(recent) != 2 {
		t.Errorf("Expected 2 recent errors, got %d", len(recent))
	}
	if recent[0] != err2 {
		t.Error("Expected most recent error to be err2")
	}
	if recent[1] != err3 {
		t.Error("Expected second most recent error to be err3")
	}
}

func TestMetrics_GetStats(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 200*time.Millisecond, true)
	m.RecordGeneration(ProviderOpenAI, "gpt-4", 50, 100*time.Millisecond, true)
	m.RecordError(ProviderOpenAI, "gpt-4", errors.New("error"))

	stats := m.GetStats(ProviderOpenAI, "gpt-4")

	if stats.Provider != ProviderOpenAI {
		t.Errorf("Expected provider %v, got %v", ProviderOpenAI, stats.Provider)
	}
	if stats.Model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", stats.Model)
	}
	if stats.TotalGenerations != 2 {
		t.Errorf("Expected 2 generations, got %d", stats.TotalGenerations)
	}
	if stats.TotalTokens != 150 {
		t.Errorf("Expected 150 tokens, got %d", stats.TotalTokens)
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics()

	m.RecordGeneration(ProviderOpenAI, "gpt-4", 100, 100*time.Millisecond, true)
	m.RecordError(ProviderOpenAI, "gpt-4", errors.New("error"))

	m.Reset()

	if len(m.generations) != 0 {
		t.Error("Generations map not cleared")
	}
	if len(m.errors) != 0 {
		t.Error("Errors map not cleared")
	}
}

func TestMetrics_key(t *testing.T) {
	m := NewMetrics()

	key1 := m.key(ProviderOpenAI, "gpt-4")
	expected1 := "openai:gpt-4"
	if key1 != expected1 {
		t.Errorf("Expected key '%s', got '%s'", expected1, key1)
	}

	key2 := m.key(ProviderOpenAI, "")
	expected2 := "openai"
	if key2 != expected2 {
		t.Errorf("Expected key '%s', got '%s'", expected2, key2)
	}
}

func TestMetrics_parseKey(t *testing.T) {
	m := NewMetrics()

	provider, model := m.parseKey("openai:gpt-4")
	if provider != ProviderOpenAI {
		t.Errorf("Expected provider %v, got %v", ProviderOpenAI, provider)
	}
	if model != "gpt-4" {
		t.Errorf("Expected model 'gpt-4', got '%s'", model)
	}

	provider2, model2 := m.parseKey("openai")
	if provider2 != ProviderOpenAI {
		t.Errorf("Expected provider %v, got %v", ProviderOpenAI, provider2)
	}
	if model2 != "" {
		t.Errorf("Expected empty model, got '%s'", model2)
	}
}
