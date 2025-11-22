package app

import (
	"context"
	"fmt"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"go.uber.org/zap"
)

// ReleaseReservationRequest represents a request to release a reservation
type ReleaseReservationRequest struct {
	ReservationID string
}

// ReleaseReservationResponse represents the response from releasing a reservation
type ReleaseReservationResponse struct {
	ReservationID string
	SKU           string
	Location      string
	Quantity      int64
	ReleasedAt    string
}

// ReleaseReservationUseCase handles reservation release
type ReleaseReservationUseCase struct {
	ledgerRepo      LedgerRepository
	reservationRepo ReservationRepository
	logger          *zap.Logger
	policy          *ledger.Policy
}

// NewReleaseReservationUseCase creates a new ReleaseReservation use case
func NewReleaseReservationUseCase(
	ledgerRepo LedgerRepository,
	reservationRepo ReservationRepository,
	logger *zap.Logger,
) *ReleaseReservationUseCase {
	return &ReleaseReservationUseCase{
		ledgerRepo:      ledgerRepo,
		reservationRepo: reservationRepo,
		logger:          logger,
		policy:          ledger.NewPolicy(),
	}
}

// Execute executes the reservation release
func (uc *ReleaseReservationUseCase) Execute(ctx context.Context, req ReleaseReservationRequest) (*ReleaseReservationResponse, error) {
	// Find reservation
	reservation, err := uc.reservationRepo.FindByID(ctx, req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("reservation not found: %w", err)
	}

	// Validate release using policy
	if err := uc.policy.ValidateRelease(reservation); err != nil {
		return nil, err
	}

	// Find ledger
	stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, reservation.SKU(), reservation.Location())
	if err != nil {
		return nil, fmt.Errorf("failed to find ledger: %w", err)
	}

	// Release reservation in ledger
	if err := stockLedger.ReleaseReservation(reservation.Quantity()); err != nil {
		return nil, err
	}

	// Release reservation entity
	if err := reservation.Release(); err != nil {
		return nil, err
	}

	// Persist changes (should be in transaction)
	if err := uc.ledgerRepo.Save(ctx, stockLedger); err != nil {
		return nil, fmt.Errorf("failed to save ledger: %w", err)
	}
	if err := uc.reservationRepo.Save(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	return &ReleaseReservationResponse{
		ReservationID: reservation.ID(),
		SKU:           reservation.SKU(),
		Location:      reservation.Location(),
		Quantity:      reservation.Quantity(),
		ReleasedAt:    reservation.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
