// Package serverless provides serverless compute implementations
package serverless

import (
	"context"
)

// AWSLambda provides AWS Lambda operations
type AWSLambda interface {
	// Invoke invokes a Lambda function
	Invoke(ctx context.Context, functionName string, payload []byte) ([]byte, error)

	// CreateFunction creates a Lambda function
	CreateFunction(ctx context.Context, config *LambdaConfig) error

	// DeleteFunction deletes a Lambda function
	DeleteFunction(ctx context.Context, functionName string) error
}

// LambdaConfig represents Lambda function configuration
type LambdaConfig struct {
	Name       string
	Runtime    string
	Handler    string
	Code       []byte
	Timeout    int
	MemorySize int
}
