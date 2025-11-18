// Package serverless provides serverless compute implementations
package serverless

import (
	"context"
)

// AzureFunctions provides Azure Functions operations
type AzureFunctions interface {
	// Invoke invokes an Azure Function
	Invoke(ctx context.Context, functionName string, payload []byte) ([]byte, error)

	// CreateFunction creates an Azure Function
	CreateFunction(ctx context.Context, config *FunctionConfig) error

	// DeleteFunction deletes an Azure Function
	DeleteFunction(ctx context.Context, functionName string) error
}

// FunctionConfig represents Azure Function configuration
type FunctionConfig struct {
	Name    string
	Runtime string
	Code    []byte
	Timeout int
}
