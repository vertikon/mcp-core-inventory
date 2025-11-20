package ledger

import (
	"errors"
	"time"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeReservation MovementType = "reservation"
	MovementTypeConfirmation MovementType = "confirmation"
	MovementTypeRelease     MovementType = "release"
	MovementTypeAdjustment   MovementType = "adjustment"
	MovementTypeReceipt     MovementType = "receipt"
	MovementTypeShipment    MovementType = "shipment"
	MovementTypeReturn      MovementType = "return"
	MovementTypeTransfer    MovementType = "transfer"
)

// StockMovement represents a movement in the stock ledger
type StockMovement struct {
	id            string
	ledgerID      string
	movementType  MovementType
	quantity      int64
	referenceID   string // Reservation ID, Order ID, etc.
	reason        string
	metadata      map[string]string
	createdAt     time.Time
	createdBy     string
}

// NewStockMovement creates a new stock movement
func NewStockMovement(id, ledgerID string, movementType MovementType, quantity int64, referenceID, reason, createdBy string) (*StockMovement, error) {
	if id == "" {
		return nil, errors.New("movement id cannot be empty")
	}
	if ledgerID == "" {
		return nil, errors.New("ledger id cannot be empty")
	}
	if movementType == "" {
		return nil, errors.New("movement type cannot be empty")
	}
	if quantity == 0 {
		return nil, errors.New("quantity cannot be zero")
	}

	return &StockMovement{
		id:           id,
		ledgerID:     ledgerID,
		movementType: movementType,
		quantity:     quantity,
		referenceID:  referenceID,
		reason:       reason,
		metadata:     make(map[string]string),
		createdAt:    time.Now(),
		createdBy:    createdBy,
	}, nil
}

// ID returns the movement ID
func (m *StockMovement) ID() string {
	return m.id
}

// LedgerID returns the ledger ID
func (m *StockMovement) LedgerID() string {
	return m.ledgerID
}

// MovementType returns the movement type
func (m *StockMovement) MovementType() MovementType {
	return m.movementType
}

// Quantity returns the movement quantity
func (m *StockMovement) Quantity() int64 {
	return m.quantity
}

// ReferenceID returns the reference ID
func (m *StockMovement) ReferenceID() string {
	return m.referenceID
}

// Reason returns the reason for the movement
func (m *StockMovement) Reason() string {
	return m.reason
}

// Metadata returns the metadata map
func (m *StockMovement) Metadata() map[string]string {
	return m.metadata
}

// SetMetadata sets a metadata key-value pair
func (m *StockMovement) SetMetadata(key, value string) {
	if m.metadata == nil {
		m.metadata = make(map[string]string)
	}
	m.metadata[key] = value
}

// CreatedAt returns the creation timestamp
func (m *StockMovement) CreatedAt() time.Time {
	return m.createdAt
}

// CreatedBy returns who created the movement
func (m *StockMovement) CreatedBy() string {
	return m.createdBy
}

