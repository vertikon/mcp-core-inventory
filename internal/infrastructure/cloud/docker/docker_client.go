// Package docker provides Docker client implementations
package docker

import (
	"context"
)

// DockerClient provides Docker API operations
type DockerClient interface {
	// Ping checks if Docker daemon is accessible
	Ping(ctx context.Context) error

	// Version returns Docker version information
	Version(ctx context.Context) (map[string]interface{}, error)

	// Close closes the Docker client connection
	Close() error
}
