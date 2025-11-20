package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"go.uber.org/zap"
)

// ReserveStockRequest represents a request to reserve stock
type ReserveStockRequest struct {
	SKU            string
	Location       string
	Quantity       int64
	IdempotencyKey string
	TTL            time.Duration
}

// ReserveStockResponse represents the response from reserving stock
type ReserveStockResponse struct {
	ReservationID   string
	IdempotencyKey string
	SKU            string
	Location       string
	Quantity       int64
	ExpiresAt      time.Time
}

// LedgerRepository defines the interface for ledger persistence
type LedgerRepository interface {
	FindBySKUAndLocation(ctx context.Context, sku, location string) (*ledger.StockLedger, error)
	Save(ctx context.Context, ledger *ledger.StockLedger) error
	SaveReservation(ctx context.Context, reservation *ledger.Reservation) error
	FindReservationByKey(ctx context.Context, idempotencyKey string) (*ledger.Reservation, error)
}

// LockService defines the interface for distributed locking
type LockService interface {
	AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	PublishReserveRequest(ctx context.Context, req ReserveStockRequest, reservationID string) error
	PublishReserveConfirmed(ctx context.Context, reservation *ledger.Reservation) error
}

// ReserveStockUseCase handles stock reservation
type ReserveStockUseCase struct {
	ledgerRepo   LedgerRepository
	lockService  LockService
	eventPub     EventPublisher
	logger       *zap.Logger
	policy       *ledger.Policy
}

// NewReserveStockUseCase creates a new ReserveStock use case
func NewReserveStockUseCase(
	ledgerRepo LedgerRepository,
	lockService LockService,
	eventPub EventPublisher,
	logger *zap.Logger,
) *ReserveStockUseCase {
	return &ReserveStockUseCase{
		ledgerRepo:  ledgerRepo,
		lockService: lockService,
		eventPub:    eventPub,
		logger:      logger,
		policy:      ledger.NewPolicy(),
	}
}

// Execute executes the stock reservation
func (uc *ReserveStockUseCase) Execute(ctx context.Context, req ReserveStockRequest) (*ReserveStockResponse, error) {
	// Generate idempotency key if not provided
	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		var err error
		idempotencyKey, err = generateIdempotencyKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate idempotency key: %w", err)
		}
	}

	// Check for existing reservation with same idempotency key (idempotency check)
	existingReservation, err := uc.ledgerRepo.FindReservationByKey(ctx, idempotencyKey)
	if err == nil && existingReservation != nil {
		// Return existing reservation
		return &ReserveStockResponse{
			ReservationID:   existingReservation.ID(),
			IdempotencyKey:  existingReservation.IdempotencyKey(),
			SKU:             existingReservation.SKU(),
			Location:        existingReservation.Location(),
			Quantity:        existingReservation.Quantity(),
			ExpiresAt:       existingReservation.ExpiresAt(),
		}, nil
	}

	// Default TTL if not provided
	ttl := req.TTL
	if ttl == 0 {
		ttl = 15 * time.Minute // Default 15 minutes
	}

	// Lock key: sku:location
	lockKey := fmt.Sprintf("reserve:%s:%s", req.SKU, req.Location)
	acquired, err := uc.lockService.AcquireLock(ctx, lockKey, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !acquired {
		return nil, fmt.Errorf("failed to acquire lock: lock already held")
	}
	defer uc.lockService.ReleaseLock(ctx, lockKey)

	// Find or create ledger
	stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, req.SKU, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to find ledger: %w", err)
	}

	// Validate reservation using policy
	if err := uc.policy.ValidateReservation(stockLedger, req.Quantity); err != nil {
		return nil, err
	}

	// Create reservation
	reservationID := generateReservationID()
	reservation, err := ledger.NewReservation(
		reservationID,
		stockLedger.ID(),
		req.SKU,
		req.Location,
		idempotencyKey,
		req.Quantity,
		ttl,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Reserve stock in ledger
	if err := stockLedger.Reserve(req.Quantity); err != nil {
		return nil, err
	}

	// Persist ledger and reservation (should be in transaction)
	if err := uc.ledgerRepo.Save(ctx, stockLedger); err != nil {
		return nil, fmt.Errorf("failed to save ledger: %w", err)
	}
	if err := uc.ledgerRepo.SaveReservation(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	// Publish event
	if err := uc.eventPub.PublishReserveRequest(ctx, req, reservationID); err != nil {
		uc.logger.Warn("failed to publish reserve request event", zap.Error(err))
		// Don't fail the operation if event publishing fails
	}

	return &ReserveStockResponse{
		ReservationID:   reservationID,
		IdempotencyKey:  idempotencyKey,
		SKU:             req.SKU,
		Location:        req.Location,
		Quantity:        req.Quantity,
		ExpiresAt:       reservation.ExpiresAt(),
	}, nil
}

func generateIdempotencyKey() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func generateReservationID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

