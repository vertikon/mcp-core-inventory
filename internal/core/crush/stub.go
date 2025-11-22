// Package crush currently provides no-op placeholders so the rest of the
// system can build while the real optimizers are being implemented.
package crush

import "context"

// OptimizerConfig represents the configuration for the optimizer.
type OptimizerConfig struct{}

// Optimizer is a no-op optimizer placeholder.
type Optimizer struct{}

// NewOptimizer returns a new optimizer stub.
func NewOptimizer(OptimizerConfig) *Optimizer {
	return &Optimizer{}
}

// Optimize does nothing and succeeds.
func (o *Optimizer) Optimize(context.Context) error {
	return nil
}

// BatchProcessorConfig represents batch-processing configuration.
type BatchProcessorConfig struct{}

// BatchProcessor is a stubbed processor that fulfils previous interfaces.
type BatchProcessor struct{}

// NewBatchProcessor returns a new BatchProcessor stub.
func NewBatchProcessor(BatchProcessorConfig) *BatchProcessor {
	return &BatchProcessor{}
}

// Start starts the processor (no-op).
func (bp *BatchProcessor) Start(context.Context) error { return nil }

// Stop stops the processor (no-op).
func (bp *BatchProcessor) Stop(context.Context) error { return nil }
