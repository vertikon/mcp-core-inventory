// Package docker provides Docker container management
package docker

import (
	"context"
)

// ContainerManager provides container management operations
type ContainerManager interface {
	// CreateContainer creates a new container
	CreateContainer(ctx context.Context, config *ContainerConfig) (string, error)

	// StartContainer starts a container
	StartContainer(ctx context.Context, containerID string) error

	// StopContainer stops a container
	StopContainer(ctx context.Context, containerID string) error

	// RemoveContainer removes a container
	RemoveContainer(ctx context.Context, containerID string) error
}

// ContainerConfig represents container configuration
type ContainerConfig struct {
	Image   string
	Command []string
	Env     map[string]string
}
