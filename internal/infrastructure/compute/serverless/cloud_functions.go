// Package serverless provides serverless compute implementations
package serverless

import (
	"context"
)

// CloudFunctions provides cloud functions operations (GCP/Azure)
type CloudFunctions interface {
	// Invoke invokes a cloud function
	Invoke(ctx context.Context, functionName string, payload []byte) ([]byte, error)

	// CreateFunction creates a cloud function
	CreateFunction(ctx context.Context, config *FunctionConfig) error

	// DeleteFunction deletes a cloud function
	DeleteFunction(ctx context.Context, functionName string) error
}

// FunctionConfig provides a minimal configuration required to create a cloud function.
type FunctionConfig struct {
	Name    string
	Runtime string
	Source  string
	Timeout int
}
