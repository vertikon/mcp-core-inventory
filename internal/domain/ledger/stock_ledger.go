// Package ledger contains the core domain entities for inventory management
package ledger

import (
	"errors"
	"time"
)

var (
	// ErrInsufficientStock indicates there is not enough stock available
	ErrInsufficientStock = errors.New("insufficient stock available")
	// ErrNegativeStock indicates an attempt to create negative stock
	ErrNegativeStock = errors.New("stock cannot be negative")
	// ErrInvalidReservation indicates an invalid reservation operation
	ErrInvalidReservation = errors.New("invalid reservation")
)

// StockLedger is the aggregate root for inventory management
// It maintains the ACID ledger of stock balances, reservations, and movements
type StockLedger struct {
	id        string
	sku       string
	location  string
	available int64 // Available stock (total - reserved)
	reserved  int64 // Reserved stock
	total     int64 // Total stock
	createdAt time.Time
	updatedAt time.Time
}

// NewStockLedger creates a new stock ledger entry
func NewStockLedger(id, sku, location string, initialQuantity int64) (*StockLedger, error) {
	if initialQuantity < 0 {
		return nil, ErrNegativeStock
	}
	if sku == "" {
		return nil, errors.New("sku cannot be empty")
	}
	if location == "" {
		return nil, errors.New("location cannot be empty")
	}

	now := time.Now()
	return &StockLedger{
		id:        id,
		sku:       sku,
		location:  location,
		available: initialQuantity,
		reserved:  0,
		total:     initialQuantity,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the ledger ID
func (l *StockLedger) ID() string {
	return l.id
}

// SKU returns the SKU
func (l *StockLedger) SKU() string {
	return l.sku
}

// Location returns the location
func (l *StockLedger) Location() string {
	return l.location
}

// Available returns the available stock quantity
func (l *StockLedger) Available() int64 {
	return l.available
}

// Reserved returns the reserved stock quantity
func (l *StockLedger) Reserved() int64 {
	return l.reserved
}

// Total returns the total stock quantity
func (l *StockLedger) Total() int64 {
	return l.total
}

// CreatedAt returns the creation timestamp
func (l *StockLedger) CreatedAt() time.Time {
	return l.createdAt
}

// UpdatedAt returns the last update timestamp
func (l *StockLedger) UpdatedAt() time.Time {
	return l.updatedAt
}

// Reserve attempts to reserve stock
// Returns error if insufficient stock available
func (l *StockLedger) Reserve(quantity int64) error {
	if quantity <= 0 {
		return errors.New("reservation quantity must be positive")
	}
	if l.available < quantity {
		return ErrInsufficientStock
	}

	l.available -= quantity
	l.reserved += quantity
	l.updatedAt = time.Now()
	return nil
}

// ConfirmReservation confirms a reservation and reduces total stock
func (l *StockLedger) ConfirmReservation(quantity int64) error {
	if quantity <= 0 {
		return errors.New("confirmation quantity must be positive")
	}
	if l.reserved < quantity {
		return ErrInvalidReservation
	}

	l.reserved -= quantity
	l.total -= quantity
	l.updatedAt = time.Now()
	return nil
}

// ReleaseReservation releases reserved stock back to available
func (l *StockLedger) ReleaseReservation(quantity int64) error {
	if quantity <= 0 {
		return errors.New("release quantity must be positive")
	}
	if l.reserved < quantity {
		return ErrInvalidReservation
	}

	l.reserved -= quantity
	l.available += quantity
	l.updatedAt = time.Now()
	return nil
}

// Adjust adjusts stock (for administrative corrections)
// Positive quantity increases stock, negative decreases
func (l *StockLedger) Adjust(quantity int64) error {
	newTotal := l.total + quantity
	if newTotal < 0 {
		return ErrNegativeStock
	}

	// Adjust available stock proportionally
	// If reducing, reduce from available first, then reserved
	if quantity < 0 {
		reduction := -quantity
		if l.available >= reduction {
			l.available -= reduction
		} else {
			remaining := reduction - l.available
			l.available = 0
			if l.reserved < remaining {
				return ErrNegativeStock
			}
			l.reserved -= remaining
		}
	} else {
		// Increase available stock
		l.available += quantity
	}

	l.total = newTotal
	l.updatedAt = time.Now()
	return nil
}

// CanReserve checks if a quantity can be reserved
func (l *StockLedger) CanReserve(quantity int64) bool {
	return l.available >= quantity && quantity > 0
}

