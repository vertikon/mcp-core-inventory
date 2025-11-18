package dtos

import "fmt"

// CreateTemplateRequest represents a request to create a template
type CreateTemplateRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Validate validates the CreateTemplateRequest
func (r *CreateTemplateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// UpdateTemplateRequest represents a request to update a template
type UpdateTemplateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Validate validates the UpdateTemplateRequest
func (r *UpdateTemplateRequest) Validate() error {
	return nil
}

// TemplateResponse represents a template response
type TemplateResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}
