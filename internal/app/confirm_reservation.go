package app

import (
	"context"
	"fmt"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"go.uber.org/zap"
)

// ConfirmReservationRequest represents a request to confirm a reservation
type ConfirmReservationRequest struct {
	ReservationID string
}

// ConfirmReservationResponse represents the response from confirming a reservation
type ConfirmReservationResponse struct {
	ReservationID string
	SKU           string
	Location      string
	Quantity      int64
	ConfirmedAt   string
}

// ReservationRepository defines the interface for reservation persistence
type ReservationRepository interface {
	FindByID(ctx context.Context, id string) (*ledger.Reservation, error)
	Save(ctx context.Context, reservation *ledger.Reservation) error
}

// ConfirmReservationUseCase handles reservation confirmation
type ConfirmReservationUseCase struct {
	ledgerRepo      LedgerRepository
	reservationRepo ReservationRepository
	eventPub        EventPublisher
	logger          *zap.Logger
	policy          *ledger.Policy
}

// NewConfirmReservationUseCase creates a new ConfirmReservation use case
func NewConfirmReservationUseCase(
	ledgerRepo LedgerRepository,
	reservationRepo ReservationRepository,
	eventPub EventPublisher,
	logger *zap.Logger,
) *ConfirmReservationUseCase {
	return &ConfirmReservationUseCase{
		ledgerRepo:      ledgerRepo,
		reservationRepo: reservationRepo,
		eventPub:        eventPub,
		logger:          logger,
		policy:          ledger.NewPolicy(),
	}
}

// Execute executes the reservation confirmation
func (uc *ConfirmReservationUseCase) Execute(ctx context.Context, req ConfirmReservationRequest) (*ConfirmReservationResponse, error) {
	// Find reservation
	reservation, err := uc.reservationRepo.FindByID(ctx, req.ReservationID)
	if err != nil {
		return nil, fmt.Errorf("reservation not found: %w", err)
	}

	// Validate confirmation using policy
	if err := uc.policy.ValidateConfirmation(reservation); err != nil {
		return nil, err
	}

	// Find ledger
	stockLedger, err := uc.ledgerRepo.FindBySKUAndLocation(ctx, reservation.SKU(), reservation.Location())
	if err != nil {
		return nil, fmt.Errorf("failed to find ledger: %w", err)
	}

	// Confirm reservation in ledger
	if err := stockLedger.ConfirmReservation(reservation.Quantity()); err != nil {
		return nil, err
	}

	// Confirm reservation entity
	if err := reservation.Confirm(); err != nil {
		return nil, err
	}

	// Persist changes (should be in transaction)
	if err := uc.ledgerRepo.Save(ctx, stockLedger); err != nil {
		return nil, fmt.Errorf("failed to save ledger: %w", err)
	}
	if err := uc.reservationRepo.Save(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to save reservation: %w", err)
	}

	// Publish event
	if err := uc.eventPub.PublishReserveConfirmed(ctx, reservation); err != nil {
		uc.logger.Warn("failed to publish reserve confirmed event", zap.Error(err))
	}

	return &ConfirmReservationResponse{
		ReservationID: reservation.ID(),
		SKU:           reservation.SKU(),
		Location:      reservation.Location(),
		Quantity:      reservation.Quantity(),
		ConfirmedAt:   reservation.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
