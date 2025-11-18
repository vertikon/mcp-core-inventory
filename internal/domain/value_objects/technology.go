// Package value_objects provides domain value objects
package value_objects

import (
	"fmt"
)

// StackType represents a valid technology stack
type StackType string

const (
	StackTypeGoPremium StackType = "go-premium"
	StackTypeTinyGo    StackType = "tinygo"
	StackTypeWeb       StackType = "web"
)

// ValidStackTypes returns all valid stack types
func ValidStackTypes() []StackType {
	return []StackType{
		StackTypeGoPremium,
		StackTypeTinyGo,
		StackTypeWeb,
	}
}

// IsValid checks if the stack type is valid
func (s StackType) IsValid() bool {
	for _, valid := range ValidStackTypes() {
		if s == valid {
			return true
		}
	}
	return false
}

// String returns the string representation
func (s StackType) String() string {
	return string(s)
}

// NewStackType creates a new StackType with validation
func NewStackType(value string) (StackType, error) {
	stack := StackType(value)
	if !stack.IsValid() {
		return "", fmt.Errorf("invalid stack type: %s. Valid types: %v", value, ValidStackTypes())
	}
	return stack, nil
}
