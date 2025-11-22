package app

import (
	"context"
	"fmt"
)

// QueryAvailableRequest represents a request to query available stock
type QueryAvailableRequest struct {
	SKU      string
	Location string
}

// QueryAvailableResponse represents the response from querying available stock
type QueryAvailableResponse struct {
	SKU       string
	Location  string
	Available int64
	Reserved  int64
	Total     int64
}

// CacheService defines the interface for caching stock availability
type CacheService interface {
	GetAvailable(ctx context.Context, sku, location string) (int64, bool)
	SetAvailable(ctx context.Context, sku, location string, quantity int64, ttl int) error
}

// QueryAvailableUseCase handles stock availability queries
type QueryAvailableUseCase struct {
	ledgerRepo   LedgerRepository
	cacheService CacheService
}

// NewQueryAvailableUseCase creates a new QueryAvailable use case
func NewQueryAvailableUseCase(
	ledgerRepo LedgerRepository,
	cacheService CacheService,
) *QueryAvailableUseCase {
	return &QueryAvailableUseCase{
		ledgerRepo:   ledgerRepo,
		cacheService: cacheService,
	}
}

// Execute executes the availability query
func (uc *QueryAvailableUseCase) Execute(ctx context.Context, req QueryAvailableRequest) (*QueryAvailableResponse, error) {
	// Try cache first
	if uc.cacheService != nil {
		if cached, found := uc.cacheService.GetAvailable(ctx, req.SKU, req.Location); found {
			// Return cached value (we don't have reserved/total in cache, so fetch from ledger)
			stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, req.SKU, req.Location)
			if err == nil {
				return &QueryAvailableResponse{
					SKU:       req.SKU,
					Location:  req.Location,
					Available: cached,
					Reserved:  stockLedger.Reserved(),
					Total:     stockLedger.Total(),
				}, nil
			}
		}
	}

	// Cache miss - fetch from ledger
	stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, req.SKU, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to find ledger: %w", err)
	}

	// Update cache
	if uc.cacheService != nil {
		uc.cacheService.SetAvailable(ctx, req.SKU, req.Location, stockLedger.Available(), 60) // 60s TTL
	}

	return &QueryAvailableResponse{
		SKU:       stockLedger.SKU(),
		Location:  stockLedger.Location(),
		Available: stockLedger.Available(),
		Reserved:  stockLedger.Reserved(),
		Total:     stockLedger.Total(),
	}, nil
}
