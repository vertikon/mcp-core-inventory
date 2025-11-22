package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/vertikon/mcp-core-inventory/internal/domain/ledger"
)

// LedgerRepository implements the ledger repository interface using PostgreSQL
type LedgerRepository struct {
	db *sql.DB
}

// NewLedgerRepository creates a new PostgreSQL ledger repository
func NewLedgerRepository(db *sql.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

// FindBySKUAndLocation finds a ledger entry by SKU and location
func (r *LedgerRepository) FindBySKUAndLocation(ctx context.Context, sku, location string) (*ledger.StockLedger, error) {
	query := `
		SELECT id, sku, location, available, reserved, total, created_at, updated_at
		FROM stock_ledger
		WHERE sku = $1 AND location = $2
		FOR UPDATE
	`

	var id, skuVal, locationVal string
	var available, reserved, total int64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, sku, location).Scan(
		&id, &skuVal, &locationVal, &available, &reserved, &total, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create new ledger entry if not found
			return r.createLedger(ctx, sku, location)
		}
		return nil, err
	}

	// Reconstruct StockLedger from DB fields
	return r.reconstructLedger(id, skuVal, locationVal, available, reserved, total, createdAt, updatedAt), nil
}

// Save saves a ledger entry
func (r *LedgerRepository) Save(ctx context.Context, l *ledger.StockLedger) error {
	query := `
		INSERT INTO stock_ledger (id, sku, location, available, reserved, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (sku, location) 
		DO UPDATE SET 
			available = EXCLUDED.available,
			reserved = EXCLUDED.reserved,
			total = EXCLUDED.total,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		l.ID(), l.SKU(), l.Location(),
		l.Available(), l.Reserved(), l.Total(),
		l.CreatedAt(), l.UpdatedAt(),
	)
	return err
}

// SaveReservation saves a reservation
func (r *LedgerRepository) SaveReservation(ctx context.Context, reservation *ledger.Reservation) error {
	query := `
		INSERT INTO reservations (id, ledger_id, sku, location, quantity, idempotency_key, status, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) 
		DO UPDATE SET 
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID(), reservation.LedgerID(), reservation.SKU(), reservation.Location(),
		reservation.Quantity(), reservation.IdempotencyKey(), string(reservation.Status()),
		reservation.ExpiresAt(), reservation.CreatedAt(), reservation.UpdatedAt(),
	)
	return err
}

// FindReservationByKey finds a reservation by idempotency key
func (r *LedgerRepository) FindReservationByKey(ctx context.Context, idempotencyKey string) (*ledger.Reservation, error) {
	query := `
		SELECT id, ledger_id, sku, location, quantity, idempotency_key, status, expires_at, created_at, updated_at
		FROM reservations
		WHERE idempotency_key = $1
	`

	var (
		id        string
		ledgerID  string
		sku       string
		location  string
		quantity  int64
		storedKey string
		status    string
		expiresAt time.Time
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, query, idempotencyKey).Scan(
		&id, &ledgerID, &sku, &location, &quantity,
		&storedKey, &status, &expiresAt, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	return ledger.NewReservationFromRecord(
		id,
		ledgerID,
		sku,
		location,
		storedKey,
		quantity,
		ledger.ReservationStatus(status),
		expiresAt,
		createdAt,
		updatedAt,
	), nil
}

// FindByID finds a reservation by ID
func (r *LedgerRepository) FindByID(ctx context.Context, id string) (*ledger.Reservation, error) {
	query := `
		SELECT id, ledger_id, sku, location, quantity, idempotency_key, status, expires_at, created_at, updated_at
		FROM reservations
		WHERE id = $1
	`

	var (
		resID       string
		ledgerID    string
		sku         string
		location    string
		quantity    int64
		idempotency string
		status      string
		expiresAt   time.Time
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&resID, &ledgerID, &sku, &location, &quantity,
		&idempotency, &status, &expiresAt, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	return ledger.NewReservationFromRecord(
		resID,
		ledgerID,
		sku,
		location,
		idempotency,
		quantity,
		ledger.ReservationStatus(status),
		expiresAt,
		createdAt,
		updatedAt,
	), nil
}

// FindExpiredReservations finds all expired reservations before the given time
func (r *LedgerRepository) FindExpiredReservations(ctx context.Context, before time.Time) ([]*ledger.Reservation, error) {
	query := `
		SELECT id, ledger_id, sku, location, quantity, idempotency_key, status, expires_at, created_at, updated_at
		FROM reservations
		WHERE status = 'pending' AND expires_at < $1
		ORDER BY expires_at ASC
		LIMIT 1000
	`

	rows, err := r.db.QueryContext(ctx, query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*ledger.Reservation
	for rows.Next() {
		var (
			id          string
			ledgerID    string
			sku         string
			location    string
			quantity    int64
			idempotency string
			status      string
			expiresAt   time.Time
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(
			&id, &ledgerID, &sku, &location, &quantity,
			&idempotency, &status, &expiresAt, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}

		reservations = append(reservations, ledger.NewReservationFromRecord(
			id,
			ledgerID,
			sku,
			location,
			idempotency,
			quantity,
			ledger.ReservationStatus(status),
			expiresAt,
			createdAt,
			updatedAt,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reservations, nil
}

// SaveLedger saves a ledger entry (alias for Save for interface compatibility)
func (r *LedgerRepository) SaveLedger(ctx context.Context, l *ledger.StockLedger) error {
	return r.Save(ctx, l)
}

// Helper methods (these would need proper implementation with reflection or builders)

func (r *LedgerRepository) createLedger(ctx context.Context, sku, location string) (*ledger.StockLedger, error) {
	id := generateID()

	ledger, err := ledger.NewStockLedger(id, sku, location, 0)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO stock_ledger (id, sku, location, available, reserved, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.ExecContext(ctx, query,
		ledger.ID(), ledger.SKU(), ledger.Location(),
		ledger.Available(), ledger.Reserved(), ledger.Total(),
		ledger.CreatedAt(), ledger.UpdatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return ledger, nil
}

func (r *LedgerRepository) reconstructLedger(id, sku, location string, available, reserved, total int64, createdAt, updatedAt time.Time) *ledger.StockLedger {
	return ledger.NewStockLedgerFromRecord(id, sku, location, available, reserved, total, createdAt, updatedAt)
}

func generateID() string {
	// Simplified ID generation - use proper UUID in production
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
