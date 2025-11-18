// Package transformer provides inference engine tests
package transformer

import (
	"context"
	"testing"
)

func TestNewInferenceEngine(t *testing.T) {
	config := GLMConfig{
		VocabSize:  10000,
		HiddenSize: 512,
		NumLayers:  2,
		NumHeads:   8,
		MaxSeqLen:  1024,
	}
	transformer := NewGLMTransformer(config)

	inferenceConfig := InferenceConfig{
		MaxTokens:     100,
		Temperature:   0.7,
		TopK:          50,
		TopP:          0.9,
		BeamWidth:     4,
		UseBeamSearch: false,
	}

	engine := NewInferenceEngine(transformer, inferenceConfig)
	if engine == nil {
		t.Fatal("NewInferenceEngine() returned nil")
	}
}

func TestInferenceEngine_Sample(t *testing.T) {
	config := GLMConfig{
		VocabSize:  10000,
		HiddenSize: 512,
		NumLayers:  2,
		NumHeads:   8,
		MaxSeqLen:  1024,
	}
	transformer := NewGLMTransformer(config)

	inferenceConfig := InferenceConfig{
		MaxTokens:     10,
		Temperature:   0.7,
		TopK:          50,
		TopP:          0.9,
		UseBeamSearch: false,
	}

	engine := NewInferenceEngine(transformer, inferenceConfig)

	input := &Tensor{
		Data:  make([]float64, 10*512),
		Shape: []int{10, 512},
	}

	ctx := context.Background()
	result, err := engine.sample(ctx, input, inferenceConfig)
	if err != nil {
		t.Fatalf("sample() error = %v", err)
	}

	if result == nil {
		t.Fatal("sample() returned nil result")
	}

	if len(result.Tokens) == 0 {
		t.Error("sample() returned empty tokens")
	}
}

func TestInferenceEngine_ApplyTemperature(t *testing.T) {
	engine := &InferenceEngine{}

	probs := []float64{0.1, 0.2, 0.3, 0.4}
	temp := 0.5

	result := engine.applyTemperature(probs, temp)

	if len(result) != len(probs) {
		t.Errorf("applyTemperature() length = %v, want %v", len(result), len(probs))
	}

	// Check normalization
	sum := 0.0
	for _, prob := range result {
		sum += prob
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("applyTemperature() sum = %v, want ~1.0", sum)
	}
}

func TestInferenceEngine_ApplyTopK(t *testing.T) {
	engine := &InferenceEngine{}

	probs := []float64{0.1, 0.2, 0.3, 0.4, 0.0}
	k := 3

	result := engine.applyTopK(probs, k)

	// Check normalization
	sum := 0.0
	for _, prob := range result {
		sum += prob
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("applyTopK() sum = %v, want ~1.0", sum)
	}
}

func TestInferenceEngine_ApplyTopP(t *testing.T) {
	engine := &InferenceEngine{}

	probs := []float64{0.1, 0.2, 0.3, 0.4}
	p := 0.7

	result := engine.applyTopP(probs, p)

	// Check normalization
	sum := 0.0
	for _, prob := range result {
		sum += prob
	}
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("applyTopP() sum = %v, want ~1.0", sum)
	}
}
