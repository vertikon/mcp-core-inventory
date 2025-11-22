package services

import (
	"context"
	"fmt"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
	"github.com/vertikon/mcp-core-inventory/internal/observability"
	"go.uber.org/zap"
)

// ReservationRepository defines the interface for reservation cleanup operations
type ReservationRepository interface {
	FindExpiredReservations(ctx context.Context, before time.Time) ([]*ledger.Reservation, error)
	SaveReservation(ctx context.Context, reservation *ledger.Reservation) error
	FindBySKUAndLocation(ctx context.Context, sku, location string) (*ledger.StockLedger, error)
	Save(ctx context.Context, ledger *ledger.StockLedger) error
}

// ReservationCleanupService handles automatic cleanup of expired reservations
type ReservationCleanupService struct {
	reservationRepo ReservationRepository
	ledgerRepo      ReservationRepository // Can reuse same interface
	logger          *zap.Logger
	metrics         *observability.InventoryMetrics
	interval        time.Duration
	batchSize       int
}

// NewReservationCleanupService creates a new cleanup service
func NewReservationCleanupService(
	reservationRepo ReservationRepository,
	ledgerRepo ReservationRepository,
	logger *zap.Logger,
	metrics *observability.InventoryMetrics,
) *ReservationCleanupService {
	return &ReservationCleanupService{
		reservationRepo: reservationRepo,
		ledgerRepo:      ledgerRepo,
		logger:          logger,
		metrics:         metrics,
		interval:        1 * time.Minute, // Run every minute
		batchSize:       100,             // Process 100 at a time
	}
}

// Start starts the cleanup service in a goroutine
func (s *ReservationCleanupService) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("Starting reservation cleanup service",
		zap.Duration("interval", s.interval),
		zap.Int("batch_size", s.batchSize))

	// Run immediately on start
	s.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping reservation cleanup service")
			return
		case <-ticker.C:
			s.cleanup(ctx)
		}
	}
}

// cleanup processes expired reservations
func (s *ReservationCleanupService) cleanup(ctx context.Context) {
	now := time.Now()

	// Find expired reservations
	expiredReservations, err := s.reservationRepo.FindExpiredReservations(ctx, now)
	if err != nil {
		s.logger.Error("Failed to find expired reservations", zap.Error(err))
		return
	}

	if len(expiredReservations) == 0 {
		return
	}

	s.logger.Info("Found expired reservations",
		zap.Int("count", len(expiredReservations)))

	// Process in batches
	for i := 0; i < len(expiredReservations); i += s.batchSize {
		end := i + s.batchSize
		if end > len(expiredReservations) {
			end = len(expiredReservations)
		}

		batch := expiredReservations[i:end]
		if err := s.processBatch(ctx, batch); err != nil {
			s.logger.Error("Failed to process batch",
				zap.Int("batch_start", i),
				zap.Int("batch_size", len(batch)),
				zap.Error(err))
			continue
		}

		s.logger.Info("Processed batch",
			zap.Int("batch_start", i),
			zap.Int("batch_size", len(batch)))
	}

	// Record metric
	if s.metrics != nil {
		for i := 0; i < len(expiredReservations); i++ {
			s.metrics.RecordReservationTTLExpired()
		}
	}
}

// processBatch processes a batch of expired reservations
func (s *ReservationCleanupService) processBatch(ctx context.Context, reservations []*ledger.Reservation) error {
	for _, reservation := range reservations {
		if err := s.processExpiredReservation(ctx, reservation); err != nil {
			s.logger.Error("Failed to process expired reservation",
				zap.String("reservation_id", reservation.ID()),
				zap.Error(err))
			// Continue processing other reservations
			continue
		}
	}
	return nil
}

// processExpiredReservation processes a single expired reservation
func (s *ReservationCleanupService) processExpiredReservation(ctx context.Context, reservation *ledger.Reservation) error {
	// Only process pending reservations
	if reservation.Status() != ledger.ReservationStatusPending {
		return nil
	}

	// Mark as expired
	reservation.Expire()

	// Release stock back to available
	stockLedger, err := s.ledgerRepo.FindBySKUAndLocation(ctx, reservation.SKU(), reservation.Location())
	if err != nil {
		return fmt.Errorf("failed to find ledger: %w", err)
	}

	// Release the reserved stock
	if err := stockLedger.ReleaseReservation(reservation.Quantity()); err != nil {
		return fmt.Errorf("failed to release reservation in ledger: %w", err)
	}

	// Save both reservation and ledger
	if err := s.reservationRepo.SaveReservation(ctx, reservation); err != nil {
		return fmt.Errorf("failed to save reservation: %w", err)
	}

	if err := s.ledgerRepo.Save(ctx, stockLedger); err != nil {
		return fmt.Errorf("failed to save ledger: %w", err)
	}

	s.logger.Debug("Processed expired reservation",
		zap.String("reservation_id", reservation.ID()),
		zap.String("sku", reservation.SKU()),
		zap.String("location", reservation.Location()),
		zap.Int64("quantity", reservation.Quantity()))

	return nil
}
