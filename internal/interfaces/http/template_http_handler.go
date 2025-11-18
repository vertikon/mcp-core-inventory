package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vertikon/mcp-hulk/internal/application/dtos"
	"github.com/vertikon/mcp-hulk/internal/services"
	"github.com/vertikon/mcp-hulk/pkg/logger"
	"go.uber.org/zap"
)

// TemplateHandler handles HTTP requests for template operations
type TemplateHandler struct {
	templateService *services.TemplateAppService
	logger          *zap.Logger
}

// NewTemplateHandler creates a new template HTTP handler
func NewTemplateHandler(templateService *services.TemplateAppService) *TemplateHandler {
	return &TemplateHandler{
		templateService: templateService,
		logger:          logger.WithContext(nil),
	}
}

// RegisterRoutes registers template routes
func (h *TemplateHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/templates", h.CreateTemplate)
	e.GET("/templates", h.ListTemplates)
	e.GET("/templates/:id", h.GetTemplate)
	e.PUT("/templates/:id", h.UpdateTemplate)
	e.DELETE("/templates/:id", h.DeleteTemplate)
}

// CreateTemplate handles POST /templates
func (h *TemplateHandler) CreateTemplate(c echo.Context) error {
	var req dtos.CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// TODO: Call service
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": "placeholder", "message": "Template creation endpoint"})
}

// ListTemplates handles GET /templates
func (h *TemplateHandler) ListTemplates(c echo.Context) error {
	// TODO: Call service
	return c.JSON(http.StatusOK, map[string]interface{}{"templates": []interface{}{}, "total": 0})
}

// GetTemplate handles GET /templates/:id
func (h *TemplateHandler) GetTemplate(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Template ID is required"})
	}
	// TODO: Call service
	return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "message": "Template retrieval endpoint"})
}

// UpdateTemplate handles PUT /templates/:id
func (h *TemplateHandler) UpdateTemplate(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Template ID is required"})
	}
	var req dtos.UpdateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	// TODO: Call service
	return c.JSON(http.StatusOK, map[string]interface{}{"id": id, "message": "Template update endpoint"})
}

// DeleteTemplate handles DELETE /templates/:id
func (h *TemplateHandler) DeleteTemplate(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Template ID is required"})
	}
	// TODO: Call service
	return c.NoContent(http.StatusNoContent)
}
