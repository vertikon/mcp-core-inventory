// Package distributed provides distributed compute implementations
package distributed

import (
	"context"
)

// DaskCluster provides Dask cluster operations
type DaskCluster interface {
	// Submit submits a task to the cluster
	Submit(ctx context.Context, task interface{}) (string, error)

	// GetResult gets the result of a task
	GetResult(ctx context.Context, taskID string) (interface{}, error)

	// Close closes the cluster connection
	Close() error
}
