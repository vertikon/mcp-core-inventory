// Package kubernetes provides Kubernetes client implementations
package kubernetes

import (
	"context"
)

// ServiceManager provides service/ingress management operations
type ServiceManager interface {
	// CreateService creates a service
	CreateService(ctx context.Context, namespace string, service *Service) error

	// GetService gets a service
	GetService(ctx context.Context, namespace string, name string) (*Service, error)

	// UpdateService updates a service
	UpdateService(ctx context.Context, namespace string, service *Service) error

	// DeleteService deletes a service
	DeleteService(ctx context.Context, namespace string, name string) error
}
