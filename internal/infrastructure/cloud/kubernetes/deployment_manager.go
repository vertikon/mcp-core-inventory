// Package kubernetes provides Kubernetes client implementations
package kubernetes

import (
	"context"
)

// DeploymentManager provides deployment management operations
type DeploymentManager interface {
	// CreateDeployment creates a deployment
	CreateDeployment(ctx context.Context, namespace string, deployment *Deployment) error

	// GetDeployment gets a deployment
	GetDeployment(ctx context.Context, namespace string, name string) (*Deployment, error)

	// UpdateDeployment updates a deployment
	UpdateDeployment(ctx context.Context, namespace string, deployment *Deployment) error

	// DeleteDeployment deletes a deployment
	DeleteDeployment(ctx context.Context, namespace string, name string) error
}
