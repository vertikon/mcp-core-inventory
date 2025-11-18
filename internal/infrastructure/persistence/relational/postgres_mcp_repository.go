// Package relational provides PostgreSQL relational database implementations
package relational

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/vertikon/mcp-hulk/internal/domain/entities"
	"github.com/vertikon/mcp-hulk/internal/domain/repositories"
	"github.com/vertikon/mcp-hulk/pkg/logger"
	"go.uber.org/zap"
)

// PostgresMCPRepository implements MCPRepository using PostgreSQL
type PostgresMCPRepository struct {
	db *sql.DB
}

// NewPostgresMCPRepository creates a new PostgreSQL MCP repository
func NewPostgresMCPRepository(db *sql.DB) repositories.MCPRepository {
	return &PostgresMCPRepository{db: db}
}

// Save saves or updates an MCP
func (r *PostgresMCPRepository) Save(ctx context.Context, mcp *entities.MCP) error {
	if mcp == nil {
		return fmt.Errorf("MCP cannot be nil")
	}

	// Serialize features
	featuresJSON, err := json.Marshal(mcp.Features())
	if err != nil {
		return fmt.Errorf("failed to marshal features: %w", err)
	}

	// Serialize context if present
	var contextJSON []byte
	if mcp.HasContext() {
		ctx := mcp.Context()
		contextData := map[string]interface{}{
			"knowledge_id": ctx.KnowledgeID(),
			"documents":    ctx.Documents(),
			"embeddings":   ctx.Embeddings(),
			"metadata":     ctx.Metadata(),
		}
		contextJSON, err = json.Marshal(contextData)
		if err != nil {
			return fmt.Errorf("failed to marshal context: %w", err)
		}
	}

	query := `
		INSERT INTO mcps (id, name, description, stack, path, features, context, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			stack = EXCLUDED.stack,
			path = EXCLUDED.path,
			features = EXCLUDED.features,
			context = EXCLUDED.context,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.ExecContext(ctx, query,
		mcp.ID(),
		mcp.Name(),
		mcp.Description(),
		string(mcp.Stack()),
		mcp.Path(),
		featuresJSON,
		contextJSON,
		mcp.CreatedAt(),
		mcp.UpdatedAt(),
	)

	if err != nil {
		logger.Error("Failed to save MCP",
			zap.String("id", mcp.ID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to save MCP: %w", err)
	}

	return nil
}

// FindByID finds an MCP by ID
func (r *PostgresMCPRepository) FindByID(ctx context.Context, id string) (*entities.MCP, error) {
	query := `
		SELECT id, name, description, stack, path, features, context, created_at, updated_at
		FROM mcps
		WHERE id = $1
	`

	var mcp entities.MCP
	var featuresJSON, contextJSON []byte
	var stack string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&mcp, // This would need proper unmarshaling
		&featuresJSON,
		&contextJSON,
		&stack,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("MCP not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find MCP: %w", err)
	}

	// TODO: Unmarshal and reconstruct entity
	// This is a placeholder - full implementation requires entity reconstruction

	return nil, fmt.Errorf("not implemented: entity reconstruction needed")
}

// FindByName finds an MCP by name
func (r *PostgresMCPRepository) FindByName(ctx context.Context, name string) (*entities.MCP, error) {
	query := `
		SELECT id, name, description, stack, path, features, context, created_at, updated_at
		FROM mcps
		WHERE name = $1
	`

	var mcp entities.MCP
	var featuresJSON, contextJSON []byte
	var stack string

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&mcp,
		&featuresJSON,
		&contextJSON,
		&stack,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("MCP not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find MCP: %w", err)
	}

	return nil, fmt.Errorf("not implemented: entity reconstruction needed")
}

// List lists all MCPs with optional filters
func (r *PostgresMCPRepository) List(ctx context.Context, filters *repositories.MCPFilters) ([]*entities.MCP, error) {
	query := "SELECT id, name, description, stack, path, features, context, created_at, updated_at FROM mcps WHERE 1=1"
	args := []interface{}{}
	argPos := 1

	if filters != nil {
		if filters.Stack != "" {
			query += fmt.Sprintf(" AND stack = $%d", argPos)
			args = append(args, filters.Stack)
			argPos++
		}
		if filters.HasContext {
			query += fmt.Sprintf(" AND context IS NOT NULL")
		}
		if filters.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argPos)
			args = append(args, filters.Limit)
			argPos++
		}
		if filters.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argPos)
			args = append(args, filters.Offset)
			argPos++
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list MCPs: %w", err)
	}
	defer rows.Close()

	var mcps []*entities.MCP
	for rows.Next() {
		// TODO: Scan and reconstruct entities
		// Placeholder implementation
	}

	return mcps, nil
}

// Delete deletes an MCP by ID
func (r *PostgresMCPRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM mcps WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete MCP: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("MCP not found: %s", id)
	}

	return nil
}

// Exists checks if an MCP exists by ID
func (r *PostgresMCPRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM mcps WHERE id = $1)"
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return exists, nil
}

// InitSchema creates the database schema if it doesn't exist
// Deprecated: Use InitAllSchemas instead
func InitSchema(db *sql.DB) error {
	return InitAllSchemas(db)
}
