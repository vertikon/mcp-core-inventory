// Package kubernetes provides Kubernetes client implementations
package kubernetes

import (
	"context"
)

// ConfigMapManager provides ConfigMap management operations
type ConfigMapManager interface {
	// CreateConfigMap creates a ConfigMap
	CreateConfigMap(ctx context.Context, namespace string, configMap *ConfigMap) error

	// GetConfigMap gets a ConfigMap
	GetConfigMap(ctx context.Context, namespace string, name string) (*ConfigMap, error)

	// UpdateConfigMap updates a ConfigMap
	UpdateConfigMap(ctx context.Context, namespace string, configMap *ConfigMap) error

	// DeleteConfigMap deletes a ConfigMap
	DeleteConfigMap(ctx context.Context, namespace string, name string) error
}
