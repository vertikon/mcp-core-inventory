package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"go.uber.org/zap"
)

// AdjustStockRequest represents a request to adjust stock
type AdjustStockRequest struct {
	SKU       string
	Location  string
	Quantity  int64 // Positive for increase, negative for decrease
	Reason    string
	CreatedBy string
}

// AdjustStockResponse represents the response from adjusting stock
type AdjustStockResponse struct {
	LedgerID     string
	SKU          string
	Location     string
	NewTotal     int64
	NewAvailable int64
}

// MovementRepository defines the interface for movement persistence
type MovementRepository interface {
	Save(ctx context.Context, movement *ledger.StockMovement) error
}

// AdjustStockUseCase handles stock adjustments
type AdjustStockUseCase struct {
	ledgerRepo   LedgerRepository
	movementRepo MovementRepository
	logger       *zap.Logger
	policy       *ledger.Policy
}

// NewAdjustStockUseCase creates a new AdjustStock use case
func NewAdjustStockUseCase(
	ledgerRepo LedgerRepository,
	movementRepo MovementRepository,
	logger *zap.Logger,
) *AdjustStockUseCase {
	return &AdjustStockUseCase{
		ledgerRepo:   ledgerRepo,
		movementRepo: movementRepo,
		logger:       logger,
		policy:       ledger.NewPolicy(),
	}
}

// Execute executes the stock adjustment
func (uc *AdjustStockUseCase) Execute(ctx context.Context, req AdjustStockRequest) (*AdjustStockResponse, error) {
	// Find or create ledger
	stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, req.SKU, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to find ledger: %w", err)
	}

	// Validate adjustment using policy
	if err := uc.policy.ValidateAdjustment(stockLedger, req.Quantity); err != nil {
		return nil, err
	}

	// Adjust stock in ledger
	if err := stockLedger.Adjust(req.Quantity); err != nil {
		return nil, err
	}

	// Create movement record
	movementID := generateMovementID()
	movement, err := ledger.NewStockMovement(
		movementID,
		stockLedger.ID(),
		ledger.MovementTypeAdjustment,
		req.Quantity,
		"", // No reference ID for adjustments
		req.Reason,
		req.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create movement: %w", err)
	}

	// Persist changes (should be in transaction)
	if err := uc.ledgerRepo.Save(ctx, stockLedger); err != nil {
		return nil, fmt.Errorf("failed to save ledger: %w", err)
	}
	if err := uc.movementRepo.Save(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to save movement: %w", err)
	}

	return &AdjustStockResponse{
		LedgerID:     stockLedger.ID(),
		SKU:          stockLedger.SKU(),
		Location:     stockLedger.Location(),
		NewTotal:     stockLedger.Total(),
		NewAvailable: stockLedger.Available(),
	}, nil
}

func generateMovementID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
