// Package value_objects provides domain value objects
package value_objects

import (
	"fmt"
)

// ValidationRuleType represents the type of validation rule
type ValidationRuleType string

const (
	ValidationRuleTypeRequired ValidationRuleType = "required"
	ValidationRuleTypeMin      ValidationRuleType = "min"
	ValidationRuleTypeMax      ValidationRuleType = "max"
	ValidationRuleTypePattern  ValidationRuleType = "pattern"
	ValidationRuleTypeCustom   ValidationRuleType = "custom"
)

// ValidationRule represents a validation rule
type ValidationRule struct {
	ruleType ValidationRuleType
	field    string
	value    interface{}
	message  string
}

// NewValidationRule creates a new validation rule
func NewValidationRule(ruleType ValidationRuleType, field string, value interface{}, message string) (*ValidationRule, error) {
	if field == "" {
		return nil, fmt.Errorf("validation rule field cannot be empty")
	}
	if message == "" {
		return nil, fmt.Errorf("validation rule message cannot be empty")
	}

	return &ValidationRule{
		ruleType: ruleType,
		field:    field,
		value:    value,
		message:  message,
	}, nil
}

// Type returns the rule type
func (v *ValidationRule) Type() ValidationRuleType {
	return v.ruleType
}

// Field returns the field name
func (v *ValidationRule) Field() string {
	return v.field
}

// Value returns the rule value
func (v *ValidationRule) Value() interface{} {
	return v.value
}

// Message returns the validation message
func (v *ValidationRule) Message() string {
	return v.message
}

// Validate validates a value against the rule
func (v *ValidationRule) Validate(value interface{}) error {
	switch v.ruleType {
	case ValidationRuleTypeRequired:
		if value == nil || value == "" {
			return fmt.Errorf(v.message)
		}
	case ValidationRuleTypeMin:
		// Type-specific validation would go here
		return nil
	case ValidationRuleTypeMax:
		// Type-specific validation would go here
		return nil
	case ValidationRuleTypePattern:
		// Pattern validation would go here
		return nil
	case ValidationRuleTypeCustom:
		// Custom validation would go here
		return nil
	default:
		return fmt.Errorf("unknown validation rule type: %s", v.ruleType)
	}
	return nil
}
