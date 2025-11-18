// Package crush provides crush optimizer tests
package crush

import (
	"context"
	"testing"
)

func TestNewOptimizer(t *testing.T) {
	config := OptimizerConfig{
		NumWorkers: 4,
		BatchSize:  16,
	}

	optimizer := NewOptimizer(config)
	if optimizer == nil {
		t.Fatal("NewOptimizer() returned nil")
	}

	if optimizer.GetWorkerCount() != config.NumWorkers {
		t.Errorf("GetWorkerCount() = %v, want %v", optimizer.GetWorkerCount(), config.NumWorkers)
	}

	if optimizer.GetBatchSize() != config.BatchSize {
		t.Errorf("GetBatchSize() = %v, want %v", optimizer.GetBatchSize(), config.BatchSize)
	}
}

func TestOptimizer_ProcessBatch(t *testing.T) {
	optimizer := NewOptimizer(OptimizerConfig{
		NumWorkers: 2,
		BatchSize:  4,
	})

	inputs := make([]interface{}, 10)
	for i := range inputs {
		inputs[i] = i
	}

	processor := func(ctx context.Context, input interface{}) (interface{}, error) {
		val := input.(int)
		return val * 2, nil
	}

	ctx := context.Background()
	results, err := optimizer.ProcessBatch(ctx, inputs, processor)
	if err != nil {
		t.Fatalf("ProcessBatch() error = %v", err)
	}

	if len(results) != len(inputs) {
		t.Errorf("ProcessBatch() results length = %v, want %v", len(results), len(inputs))
	}

	for i, result := range results {
		expected := inputs[i].(int) * 2
		if result.(int) != expected {
			t.Errorf("ProcessBatch() result[%d] = %v, want %v", i, result, expected)
		}
	}
}

func TestOptimizer_CreateBatches(t *testing.T) {
	optimizer := NewOptimizer(OptimizerConfig{
		BatchSize: 3,
	})

	inputs := make([]interface{}, 10)
	for i := range inputs {
		inputs[i] = i
	}

	batches := optimizer.createBatches(inputs, 3)
	expectedBatches := 4 // ceil(10/3)

	if len(batches) != expectedBatches {
		t.Errorf("createBatches() batches = %v, want %v", len(batches), expectedBatches)
	}
}
