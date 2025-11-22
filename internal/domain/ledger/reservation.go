package ledger

import (
	"errors"
	"time"
)

var (
	// ErrReservationExpired indicates the reservation has expired
	ErrReservationExpired = errors.New("reservation expired")
	// ErrReservationNotFound indicates the reservation does not exist
	ErrReservationNotFound = errors.New("reservation not found")
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusReleased  ReservationStatus = "released"
	ReservationStatusExpired   ReservationStatus = "expired"
)

// Reservation represents a stock reservation
type Reservation struct {
	id             string
	ledgerID       string
	sku            string
	location       string
	quantity       int64
	idempotencyKey string
	status         ReservationStatus
	expiresAt      time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

// NewReservation creates a new reservation
func NewReservation(id, ledgerID, sku, location, idempotencyKey string, quantity int64, ttl time.Duration) (*Reservation, error) {
	if id == "" {
		return nil, errors.New("reservation id cannot be empty")
	}
	if ledgerID == "" {
		return nil, errors.New("ledger id cannot be empty")
	}
	if sku == "" {
		return nil, errors.New("sku cannot be empty")
	}
	if location == "" {
		return nil, errors.New("location cannot be empty")
	}
	if idempotencyKey == "" {
		return nil, errors.New("idempotency key cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	if ttl <= 0 {
		return nil, errors.New("ttl must be positive")
	}

	now := time.Now()
	return &Reservation{
		id:             id,
		ledgerID:       ledgerID,
		sku:            sku,
		location:       location,
		quantity:       quantity,
		idempotencyKey: idempotencyKey,
		status:         ReservationStatusPending,
		expiresAt:      now.Add(ttl),
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ID returns the reservation ID
func (r *Reservation) ID() string {
	return r.id
}

// LedgerID returns the ledger ID
func (r *Reservation) LedgerID() string {
	return r.ledgerID
}

// SKU returns the SKU
func (r *Reservation) SKU() string {
	return r.sku
}

// Location returns the location
func (r *Reservation) Location() string {
	return r.location
}

// Quantity returns the reserved quantity
func (r *Reservation) Quantity() int64 {
	return r.quantity
}

// IdempotencyKey returns the idempotency key
func (r *Reservation) IdempotencyKey() string {
	return r.idempotencyKey
}

// Status returns the reservation status
func (r *Reservation) Status() ReservationStatus {
	return r.status
}

// ExpiresAt returns the expiration time
func (r *Reservation) ExpiresAt() time.Time {
	return r.expiresAt
}

// CreatedAt returns the creation timestamp
func (r *Reservation) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the last update timestamp
func (r *Reservation) UpdatedAt() time.Time {
	return r.updatedAt
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.expiresAt)
}

// Confirm confirms the reservation
func (r *Reservation) Confirm() error {
	if r.status != ReservationStatusPending {
		return errors.New("only pending reservations can be confirmed")
	}
	if r.IsExpired() {
		return ErrReservationExpired
	}

	r.status = ReservationStatusConfirmed
	r.updatedAt = time.Now()
	return nil
}

// Release releases the reservation
func (r *Reservation) Release() error {
	if r.status != ReservationStatusPending {
		return errors.New("only pending reservations can be released")
	}

	r.status = ReservationStatusReleased
	r.updatedAt = time.Now()
	return nil
}

// Expire marks the reservation as expired
func (r *Reservation) Expire() {
	if r.status == ReservationStatusPending {
		r.status = ReservationStatusExpired
		r.updatedAt = time.Now()
	}
}

// NewReservationFromRecord reconstructs a reservation from persisted data.
func NewReservationFromRecord(
	id, ledgerID, sku, location, idempotencyKey string,
	quantity int64,
	status ReservationStatus,
	expiresAt, createdAt, updatedAt time.Time,
) *Reservation {
	return &Reservation{
		id:             id,
		ledgerID:       ledgerID,
		sku:            sku,
		location:       location,
		quantity:       quantity,
		idempotencyKey: idempotencyKey,
		status:         status,
		expiresAt:      expiresAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}
