// Package crush provides Crush optimizations for GLM-4.6
package crush

import (
	"context"
	"runtime"
	"sync"
)

// Optimizer provides Crush optimization capabilities
type Optimizer struct {
	numWorkers int
	batchSize  int
}

// OptimizerConfig represents optimizer configuration
type OptimizerConfig struct {
	NumWorkers int `json:"num_workers"`
	BatchSize  int `json:"batch_size"`
}

// NewOptimizer creates a new Crush optimizer
func NewOptimizer(config OptimizerConfig) *Optimizer {
	numWorkers := config.NumWorkers
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU() * 2
	}

	batchSize := config.BatchSize
	if batchSize == 0 {
		batchSize = 32
	}

	return &Optimizer{
		numWorkers: numWorkers,
		batchSize:  batchSize,
	}
}

// ProcessBatch processes a batch of inputs in parallel
func (o *Optimizer) ProcessBatch(ctx context.Context, inputs []interface{}, processor func(context.Context, interface{}) (interface{}, error)) ([]interface{}, error) {
	if len(inputs) == 0 {
		return []interface{}{}, nil
	}

	// Create batches
	batches := o.createBatches(inputs, o.batchSize)
	results := make([]interface{}, len(inputs))
	errors := make([]error, len(batches))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, o.numWorkers)

	for i, batch := range batches {
		wg.Add(1)
		go func(batchIdx int, batch []interface{}) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				errors[batchIdx] = ctx.Err()
				return
			default:
			}

			for j, input := range batch {
				result, err := processor(ctx, input)
				if err != nil {
					errors[batchIdx] = err
					return
				}
				results[batchIdx*o.batchSize+j] = result
			}
		}(i, batch)
	}

	wg.Wait()

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// createBatches splits inputs into batches
func (o *Optimizer) createBatches(inputs []interface{}, batchSize int) [][]interface{} {
	batches := make([][]interface{}, 0)
	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		batches = append(batches, inputs[i:end])
	}
	return batches
}

// CompactMemory performs memory compaction
func (o *Optimizer) CompactMemory() {
	runtime.GC()
}

// GetWorkerCount returns the number of workers
func (o *Optimizer) GetWorkerCount() int {
	return o.numWorkers
}

// GetBatchSize returns the batch size
func (o *Optimizer) GetBatchSize() int {
	return o.batchSize
}
